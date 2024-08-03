package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"context"
	"log"

	"github.com/babu10103/mongo-golang/controllers"
	"github.com/babu10103/mongo-golang/routes"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	port := os.Getenv("APP_CONTAINER_PORT")
	if port == "" {
		log.Println("APP_CONTAINER_PORT not set, defaulting to 8080")
		port = "8080"
	}
	client, err := getSession()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())
	uc := controllers.NewUserController(client)
	r := mux.NewRouter()

	routes.RegisterUserRoutes(r, *uc)

	http.ListenAndServe(":"+port, r)
}

func getSession() (*mongo.Client, error) {
	dbHost := os.Getenv("MONGODB_HOSTNAME")
	dbUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	dbPwd := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")

	if dbHost == "" || dbUser == "" || dbPwd == "" {
		return nil, fmt.Errorf("invalid dbHost, dbUser or dbPwd provided")
	}

	connStr := fmt.Sprintf("mongodb://%s:%s@%s:27017", dbUser, dbPwd, dbHost)
	clientOptions := options.Client().ApplyURI(connStr)

	var client *mongo.Client
	var err error

	for i := 0; i < 10; i++ {
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err == nil {
			// check if connection is functional and able to communicate with MongoDB server
			err = client.Ping(context.Background(), nil)
			if err == nil {
				return client, nil
			}
		}
		log.Printf("Attempt %d: Error connecting to MongoDB: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to MongoDB after multiple attempts: %w")
}
