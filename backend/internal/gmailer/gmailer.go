package gmailer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// GmailClient encapsulates Gmail API service and templates
type GmailClient struct {
	service     *gmail.Service
	templateDir string
}

// NewGmailClient initializes a Gmail client
func NewGmailClient(credentialsPath, tokenPath, templateDir string) (*GmailClient, error) {
	// Load credentials
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}

	// Get token
	tok, err := getToken(tokenPath, config)
	if err != nil {
		return nil, err
	}

	// Use TokenSource for auto-refresh
	tokenSource := config.TokenSource(context.Background(), tok)
	client := oauth2.NewClient(context.Background(), tokenSource)

	srv, err := gmail.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to create Gmail service: %v", err)
	}

	return &GmailClient{
		service:     srv,
		templateDir: templateDir,
	}, nil
}

// --- Public Email Methods ---

func (g *GmailClient) SendOTPEmail(to string, data TemplateData) error {
	body, err := LoadTemplate(g.templateDir, "otp", data)
	if err != nil {
		return err
	}
	subject := "Your OTP Verification Code"
	return g.sendEmail(to, subject, body)
}

func (g *GmailClient) SendWelcomeEmail(to string, data TemplateData) error {
	body, err := LoadTemplate(g.templateDir, "welcome", data)
	if err != nil {
		return err
	}
	subject := fmt.Sprintf("Welcome, %s!", data["Name"])
	return g.sendEmail(to, subject, body)
}

// Generic send
func (g *GmailClient) sendEmail(to, subject, body string) error {
	raw := fmt.Sprintf(
		"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		to, subject, body,
	)
	msg := &gmail.Message{Raw: encodeWeb64String([]byte(raw))}

	_, err := g.service.Users.Messages.Send("me", msg).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}
	return nil
}

// --- Helpers ---

func encodeWeb64String(b []byte) string {
	s := make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(s, b)
	return string(s)
}

// --- Token Management ---

func getToken(path string, config *oauth2.Config) (*oauth2.Token, error) {
	// Try to load saved token
	tok, err := loadToken(path)
	if err == nil {
		return tok, nil
	}

	// Fetch new token from web
	tok = getTokenFromWeb(config)
	if err := saveToken(path, tok); err != nil {
		return nil, fmt.Errorf("unable to save token: %v", err)
	}
	return tok, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		panic(err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		panic(fmt.Errorf("unable to retrieve token from web: %v", err))
	}
	return tok
}

func loadToken(path string) (*oauth2.Token, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tok oauth2.Token
	if err := json.Unmarshal(b, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func saveToken(path string, token *oauth2.Token) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
