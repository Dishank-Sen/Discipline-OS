package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Dishank-Sen/Discipline-OS/cmd/api"
	"github.com/Dishank-Sen/Discipline-OS/db/connect"
	"github.com/Dishank-Sen/Discipline-OS/internal/gmailer"
	errorhandler "github.com/Dishank-Sen/Discipline-OS/utils/errorHandler"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	errorhandler.HandleError(err, "Failed to Load env")
	port := os.Getenv("PORT")
	mongodb_uri := os.Getenv("MONGODB_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := connect.NewMongoDBStorage(mongodb_uri, ctx)
	errorhandler.HandleError(err, "Connecting to MongoDB")

	defer func() {
        if err := client.Disconnect(ctx); err != nil {
            errorhandler.HandleError(err, "Disconnecting MongoDB")
        }
    }()

	err = client.Ping(ctx, nil)
	errorhandler.HandleError(err, "Pinging MongoDB")
	fmt.Println("Connected to MongoDB")
	
	// new gmail client
	credentials := "./assets/credentials.json"
	tokenFile := "./assets/credentials-gmail.json"
	templateDir := "./internal/templates"

	gmailClient, err := gmailer.NewGmailClient(credentials, tokenFile, templateDir)
	if err != nil {
		log.Fatalf("Failed to create Gmail client: %v", err)
	}


	server := api.NewAPIServer(port, client, gmailClient)
	err = server.Run()
	errorhandler.HandleError(err, "Failed to start server")
}