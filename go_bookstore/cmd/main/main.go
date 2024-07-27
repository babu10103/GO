package main

import (
	"net/http"
	"os"

	"github.com/babu10103/GO/go_bookstore/pkg/log"
	"github.com/babu10103/GO/go_bookstore/pkg/routes"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	APP_PORT := os.Getenv("APP_PORT")
	if APP_PORT == "" {
		APP_PORT = "8080"
	}
	r := mux.NewRouter()

	routes.RegisterBookStoreRoutes(r)

	http.Handle("/", r)
	log.ErrorLogger.Fatal(http.ListenAndServe(":"+APP_PORT, r))
}
