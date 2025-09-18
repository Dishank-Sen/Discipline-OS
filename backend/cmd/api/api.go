package api

import (
	"log"
	"net/http"
	"os"

	"github.com/Dishank-Sen/Discipline-OS/service/routes"
	"github.com/Dishank-Sen/Discipline-OS/service/store"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct{
	addr string
	client *mongo.Client
}

func NewAPIServer(addr string, client *mongo.Client) *APIServer{
	return &APIServer{
		addr: addr,
		client: client,
	}
}

func (s *APIServer) Run() error{
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	db_name := os.Getenv("DB_NAME")
	userCollection_name := os.Getenv("USER_COLLECTION")
	tempUserCollection_name := os.Getenv("TEMP_USER_COLLECTION")
	userCollection := s.client.Database(db_name).Collection(userCollection_name)
	tempUserCollection := s.client.Database(db_name).Collection(tempUserCollection_name)
	userStore := store.NewStore(s.client, userCollection, tempUserCollection)
	handler := routes.NewHandler(userStore)
	handler.RegisterRoutes(subRouter)

	log.Printf("server running on port %s",s.addr)
	return http.ListenAndServe(s.addr, router)
}