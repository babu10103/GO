package main

import (
	"github.com/babu10103/GO/go_bookstore/pkg/routes"
	"github.com/gorilla/mux"
	"github.com/junzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	routes.RegisterBookStoreRoutes(r)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9010", r))

}
