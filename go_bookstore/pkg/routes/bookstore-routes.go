package routes

import (
	"github.com/babu10103/GO/go_bookstore/pkg/controllers"
	"github.com/gorilla/mux"
)

var RegisterBookStoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/books/", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/books/", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/books/{bookId}", controllers.GetBookById).Methods("GET")
	router.HandleFunc("/books/{bookId}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{bookId}", controllers.DeleteBook).Methods("DELETE")
}
