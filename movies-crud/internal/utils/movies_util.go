package utils

import (
	"math/rand"
	"movies-crud/internal/models"
	"movies-crud/pkg/log"
	"strconv"
	"sync"
	"time"
)

var (
	idCounter int
	mu        sync.Mutex
)

func init() {
	/*
		the default behavior of Go's math/rand package is to use a
		deterministic seed (e.g., 1), which results in the same sequence
		of pseudo-random numbers each time the program is run
	*/
	rand.Seed(time.Now().UnixNano())
}

func GenerateID() string {
	// Lock the mutex so that only one goroutine at a time can access it
	// Without proper synchronization, you could have race conditions
	// where two goroutines end up with the same ID or other unexpected behavior.
	mu.Lock()
	defer mu.Unlock()
	idCounter++
	return strconv.Itoa(idCounter)
}

func GenerateUID() string {
	// Define the Character Set
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	// Create a Slice to Hold the Random Characters
	b := make([]rune, 10)

	// Generate Random Characters
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetIndex(movies []models.Movie, movieID int) int {
	index := -1
	start := 0
	end := len(movies) - 1

	for start <= end {
		mid := start + (end-start)/2
		curID, err := strconv.Atoi(movies[mid].ID)
		if err != nil {
			log.ErrorLogger.Printf("Invalid movie ID in storage: %v", movies[mid].ID)
			return index
		}

		if curID == movieID {
			index = mid
			break
		} else if curID > movieID {
			end = mid - 1
		} else {
			start = mid + 1
		}
	}
	return index
}
