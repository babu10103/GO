package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"movies-crud/internal/models"
	"movies-crud/internal/utils"
	"movies-crud/pkg/log"
	"net/http"
	"strconv"
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
		log.ErrorLogger.Printf("Unsupported HTTP method: %v", r.Method)
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// newMovie := new(Movie)
	var newMovie models.Movie

	if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
		log.ErrorLogger.Printf("Failed to decode JSON body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if newMovie.Title == "" || newMovie.Director == nil || newMovie.Director.FirstName == "" || newMovie.Director.LastName == "" {
		log.ErrorLogger.Printf("Invalid movie data: %+v", newMovie)
		http.Error(w, "Invalid movie data", http.StatusBadRequest)
		return
	}

	newMovie.ID = utils.GenerateID()
	newMovie.UID = utils.GenerateUID()
	movies = append(movies, newMovie)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newMovie); err != nil {
		log.ErrorLogger.Printf("Failed to encode JSON body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.InfoLogger.Println("Request URL:", r.URL.Path)

	if r.Method != http.MethodDelete {
		log.ErrorLogger.Printf("Unsupported HTTP method: %v", r.Method)
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	movieIDStr, ok := vars["id"]
	if !ok {
		log.ErrorLogger.Println("ID not found in URL")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		log.ErrorLogger.Printf("Invalid movie ID: %v", movieIDStr)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	index := utils.GetIndex(movies, movieID)

	if index == -1 {
		log.ErrorLogger.Printf("Movie not found: %v", movieID)
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	deletedMovie := movies[index]
	movies = append(movies[:index], movies[index+1:]...)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(deletedMovie); err != nil {
		log.ErrorLogger.Printf("Failed to encode JSON body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.InfoLogger.Println("Request URL:", r.URL.Path)

	if r.Method != http.MethodGet {
		log.ErrorLogger.Printf("Unsupported HTTP method: %v", r.Method)
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	movieIDStr, ok := vars["id"]
	if !ok {
		log.ErrorLogger.Println("ID not found in URL")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		log.ErrorLogger.Printf("Invalid movie ID: %v", movieIDStr)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	index := utils.GetIndex(movies, movieID)
	if index == -1 {
		log.ErrorLogger.Printf("Movie not found: %v", movieID)
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	movie := movies[index]
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.ErrorLogger.Printf("Failed to encode JSON body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
