package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/babu10103/GO/go_bookstore/pkg/log"
	"github.com/babu10103/GO/go_bookstore/pkg/models"
	"github.com/babu10103/GO/go_bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

var NewBook models.Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
	newBooks, err := models.GetAllBooks()
	if err != nil {
		log.ErrorLogger.Fatalf("Error while getting all books:%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(newBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	ID, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		log.ErrorLogger.Fatalf("Error while parsing:%v", err)
	}
	bookDetails, _ := models.GetBookById(ID)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	utils.ParseBody(r, book)
	b := book.CreateBook()
	res, _ := json.Marshal(b)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]

	id, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		log.ErrorLogger.Fatal("Error While Parsing")
	}

	deletedBook := models.DeleteBook(id)
	log.InfoLogger.Printf("Book Deleted Successfully. name: %s, author: %s, publication: %s\n", deletedBook.Name, deletedBook.Author, deletedBook.Publication)
	res, _ := json.Marshal(deletedBook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook = &models.Book{}
	utils.ParseBody(r, updateBook)
	vars := mux.Vars(r)
	bookId := vars["bookId"]

	id, err := strconv.ParseInt(bookId, 0, 0)

	if err != nil {
		log.ErrorLogger.Fatalf("Error Whole Parsing")
	}
	bookDetails, db := models.GetBookById(id)

	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Publication
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}

	db.Save(&bookDetails)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(res)
}
