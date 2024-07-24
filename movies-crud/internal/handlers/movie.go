package handlers

import (
	"encoding/json"
	"movies-crud/internal/models"
	"movies-crud/internal/utils"
	"movies-crud/pkg/log"
	"net/http"
)

var movies = []models.Movie{}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.ErrorLogger.Printf("Failed to encode JSON body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func AddMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// newMovie := new(Movie)
	var newMovie models.Movie

	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		log.ErrorLogger.Printf("Failed to encode JSON body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if newMovie.Title == "" || newMovie.Director == nil || newMovie.Director.FirstName == "" || newMovie.Director.LastName == "" {
		http.Error(w, "Invalind movie data", http.StatusBadRequest)
		return
	}

	newMovie.ID = utils.GenerateID()
	newMovie.UID = utils.GenerateUID()
	movies = append(movies, newMovie)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newMovie); err != nil {
		log.ErrorLogger.Printf("Failed to decode JSON body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
