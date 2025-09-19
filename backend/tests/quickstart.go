package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)


// saveToken writes the OAuth token to a file
func saveToken(path string, token *oauth2.Token) {
	if token == nil {
		log.Fatal("token is nil, cannot save")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalf("Unable to create directory for token file: %v", err)
	}

	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create token file: %v", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(token); err != nil {
		log.Fatalf("Unable to encode token to file: %v", err)
	}

	log.Printf("Token saved to %s\n", path)
}

// getTokenFromWeb fetches a token via OAuth flow
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	var authCode string
	fmt.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatal(err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// getClient returns an HTTP client using a saved or new token
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "./assets/credentials-gmail.json"
	var tok *oauth2.Token

	f, err := os.Open(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	} else {
		defer f.Close()
		if err := json.NewDecoder(f).Decode(&tok); err != nil || tok == nil {
			log.Println("Failed to decode token, fetching a new one...")
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
		}
	}

	return config.Client(context.Background(), tok)
}

// encodeWeb64String encodes a byte array to base64 URL encoding without padding
func encodeWeb64String(b []byte) string {
	s := make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(s, b)
	return string(s)
}

func main() {
	// Load Google credentials
	b, err := os.ReadFile("./assets/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client credentials: %v", err)
	}

	client := getClient(config)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Prepare email
	emailTo := "dishanksen05@gmail.com"
	subject := "Test Email from Go"
	body := "This is a test email sent via Gmail API in Go 🚀"

	raw := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", emailTo, subject, body)
	message := &gmail.Message{
		Raw: encodeWeb64String([]byte(raw)),
	}

	// Send email
	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		log.Fatalf("Unable to send email: %v", err)
	}

	fmt.Println("Email sent successfully!")
}
