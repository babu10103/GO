package router

import (
	"github.com/gorilla/mux"
	"movies-crud/internal/handlers"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/movies", handlers.GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", handlers.GetMovie).Methods("GET")
	r.HandleFunc("/movies/add", handlers.AddMovie).Methods("POST")
	r.HandleFunc(
		"/movies/{id}/delete",
		handlers.DeleteMovie,
	).Methods("DELETE")
	return r
}
