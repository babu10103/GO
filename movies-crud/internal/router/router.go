package router

import (
	"movies-crud/internal/handlers"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/movies", handlers.GetMovies)
	mux.HandleFunc("/movies/add", handlers.AddMovie)
	return mux
}
