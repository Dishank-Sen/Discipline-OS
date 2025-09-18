package store

import (
	"context"
	"fmt"
	"time"

	"github.com/Dishank-Sen/Discipline-OS/service/auth"
	types "github.com/Dishank-Sen/Discipline-OS/types/database"
	errorhandler "github.com/Dishank-Sen/Discipline-OS/utils/errorHandler"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct{
	Client *mongo.Client
	UserCollection *mongo.Collection
	TempUserCollection *mongo.Collection
}

func NewStore(client *mongo.Client, userCollection *mongo.Collection, tempUserCollection *mongo.Collection) *Store{
	return &Store{
		Client: client,
		UserCollection: userCollection,
		TempUserCollection: tempUserCollection,
	}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error){
	user := &types.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	err := s.UserCollection.FindOne(ctx, filter).Decode(user)
	if err != nil{
		if err == mongo.ErrNoDocuments {
			// user not exist
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *Store) InsertEmail(email string) (string, error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existingUser, err := s.GetUserByEmail(email)
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

	_, err = s.TempUserCollection.InsertOne(ctx, data)
	return signupToken, err
}

func (s *Store) InsertPassword(password string, signupToken string) (error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hashedPassword, err := auth.HashPassword(password)
	errorhandler.HandleError(err, "Hashing password")

	filter := bson.M{"signupToken": signupToken}
	update := bson.M{"$set": bson.M{
		"password":  hashedPassword,
		"updatedAt": time.Now().UTC(),
	}}

	res, err := s.TempUserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no pending signup found for token: %s", signupToken)
	}

	return nil
}

func (s *Store) InsertOTP(otp int, signupToken string) (error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"signupToken": signupToken}
	update := bson.M{"$set": bson.M{
		"otp":  otp,
		"updatedAt": time.Now().UTC(),
	}}

	res, err := s.TempUserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no pending signup found for token: %s", signupToken)
	}

	return nil
}

func (s *Store) VerifyOTP(otp int, signupToken string) (bool){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"signupToken": signupToken,
		"otp": otp,
	}

	res := s.TempUserCollection.FindOne(ctx, filter)

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


func (s *Store) CreateNewUser(payload types.User) (string, error){
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	existingUser, err := s.GetUserByEmail(payload.Email)

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

	res, err := s.UserCollection.InsertOne(ctx, data)

	errorhandler.HandleError(err, "Creating New User")

	id, ok := res.InsertedID.(primitive.ObjectID)
	
	if !ok {
		return "", fmt.Errorf("inserted ID is not an ObjectID")
	}

	return id.Hex(), nil
}