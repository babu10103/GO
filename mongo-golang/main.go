package main

import (
	"fmt"
	"net/http"

	"context"
	"log"

	"github.com/babu10103/mongo-golang/controllers"
	"github.com/babu10103/mongo-golang/routes"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := getSession()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(context.Background())
	uc := controllers.NewUserController(client)
	r := mux.NewRouter()

	routes.RegisterUserRoutes(r, *uc)

	http.ListenAndServe(":9000", r)
}

func getSession() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)

	/*
		%w -> 	it is used with fmt.Errorf() function to wrap errors.
				wrapping errors allows you to create a new error that
				retains original error's details
	*/
	if err != nil {
		return nil, fmt.Errorf("Error connecting to MongoDB: %w", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}
