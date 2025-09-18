package connect

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBStorage(mongodb_uri string, ctx context.Context) (*mongo.Client, error){
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb_uri))
	return client, err
}