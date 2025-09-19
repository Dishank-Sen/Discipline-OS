package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Dishank-Sen/Discipline-OS/internal/gmailer"
	"github.com/Dishank-Sen/Discipline-OS/service/auth"
	types "github.com/Dishank-Sen/Discipline-OS/types/database"
	"github.com/Dishank-Sen/Discipline-OS/utils"
	errorhandler "github.com/Dishank-Sen/Discipline-OS/utils/errorHandler"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct{
	Client *mongo.Client
	GmailClient *gmailer.GmailClient
}

func NewStore(client *mongo.Client, gmailClient *gmailer.GmailClient) *Store{
	return &Store{
		Client: client,
		GmailClient: gmailClient,
	}
}

func (s *Store) DeleteRecord(ctx context.Context, collection *mongo.Collection, filter bson.M) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (s *Store) GetUserByEmail(email string, collection *mongo.Collection) (*types.User, error){
	user := &types.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil{
		if err == mongo.ErrNoDocuments {
			// user not exist
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *Store) InsertEmail(email string, collection *mongo.Collection) (string, error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existingUser, err := s.GetUserByEmail(email, collection)
	fmt.Println("existing user:", existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err // actual error
	}
	if existingUser != nil {
		return "", fmt.Errorf("user already exists")
	}


	// create a signup token
	signupToken := uuid.New().String()

	// data to insert
	data := bson.D{
		{Key: "email", Value: email},
		{Key: "signupToken", Value: signupToken},
		{Key: "updatedAt", Value: time.Now()},
	}

	_, err = collection.InsertOne(ctx, data)
	return signupToken, err
}

func (s *Store) InsertPassword(password string, signupToken string, collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	errorhandler.HandleError(err, "Hashing password")

	// Update password in temp collection
	filter := bson.M{"signupToken": signupToken}
	update := bson.M{"$set": bson.M{
		"password":  hashedPassword,
		"updatedAt": time.Now().UTC(),
	}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no pending signup found for token: %s", signupToken)
	}

	// Generate OTP
	otp, err := utils.GenerateOTPInt()
	if err != nil {
		return fmt.Errorf("error generating OTP: %v", err)
	}

	// Insert OTP to temp collection
	err = s.InsertOTP(otp, signupToken, collection)
	if err != nil {
		return fmt.Errorf("error inserting OTP: %v", err)
	}

	var userEmail string

	// Send OTP email
	data := gmailer.TemplateData{
		"OTP":           otp,
		"ExpiryMinutes": 10,
	}

	env := os.Getenv("ENV")
	if env == "development"{
		userEmail = "dishanksen05@gmail.com"
	}else{
		// Fetch user's email
		var tempUser struct {
			Email string `bson:"email"`
		}
		err = collection.FindOne(ctx, bson.M{"signupToken": signupToken}).Decode(&tempUser)
		if err != nil {
			return fmt.Errorf("failed to find user email: %v", err)
		}
		userEmail = tempUser.Email
	}

	err = s.GmailClient.SendOTPEmail(userEmail, data)
	if err != nil {
		return fmt.Errorf("error sending OTP email: %v", err)
	}

	return nil
}


func (s *Store) InsertOTP(otp int, signupToken string, collection *mongo.Collection) (error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"signupToken": signupToken}
	update := bson.M{"$set": bson.M{
		"otp":  otp,
		"updatedAt": time.Now().UTC(),
	}}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no pending signup found for token: %s", signupToken)
	}

	return nil
}

func (s *Store) VerifyOTP(otp int, signupToken string, collection *mongo.Collection) (bool){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"signupToken": signupToken,
		"otp": otp,
	}

	res := collection.FindOne(ctx, filter)

	if res.Err() != nil {
		// no user
		if res.Err() == mongo.ErrNoDocuments {
			return false
		}
		fmt.Println(res.Err())
		return false
	}

	return true
}


func (s *Store) CreateNewUser(payload types.User, collection *mongo.Collection) (string, error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existingUser, err := s.GetUserByEmail(payload.Email, collection)

	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}

	if existingUser != nil {
		return "", fmt.Errorf("user already exists")
	}

	// data to insert
	data := bson.M{
		"username": payload.UserName,
		"email": payload.Email,
		"password": payload.Password,
		"subscribed": payload.Subscribed,
		"personal": payload.Personal,
		"device": payload.Device,
		"createdAt": time.Now(),
	}

	res, err := collection.InsertOne(ctx, data)

	errorhandler.HandleError(err, "Creating New User")

	id, ok := res.InsertedID.(primitive.ObjectID)
	
	if !ok {
		return "", fmt.Errorf("inserted ID is not an ObjectID")
	}

	return id.Hex(), nil
}