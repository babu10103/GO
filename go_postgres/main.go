package main

import (
	"log"
	"net/http"
	"os"

	"github.com/babu10103/go/go-postgres/db"
	"github.com/babu10103/go/go-postgres/router"
	"github.com/gorilla/mux"
)

// GetDB returns the database connection

func main() {
	r := mux.NewRouter()
	router.StocksRoutes(r)
	if err := db.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.GetDB().Close()
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Println("APP_PORT not set, defaulting to 8080")
		port = "8080"
	}
	http.ListenAndServe(":"+port, r)
}
