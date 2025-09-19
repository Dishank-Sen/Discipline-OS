package interfaces

import (
	"context"

	types "github.com/Dishank-Sen/Discipline-OS/types/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface{
	GetUserByEmail(email string, collection *mongo.Collection) (*types.User, error)
	InsertEmail(email string, collection *mongo.Collection) (string, error)
	InsertPassword(password string, signupToken string, collection *mongo.Collection) (error)
	InsertOTP(otp int, signupToken string, collection *mongo.Collection) (error)
	VerifyOTP(otp int, signupToken string, collection *mongo.Collection) (bool)
	CreateNewUser(payload types.User, collection *mongo.Collection) (string, error)
	DeleteRecord(ctx context.Context, collection *mongo.Collection, filter bson.M) (int64, error)
}