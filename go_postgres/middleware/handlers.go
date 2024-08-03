package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/babu10103/go/go-postgres/db"
	"github.com/babu10103/go/go-postgres/models"
	"github.com/gorilla/mux"
)

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// GetStock handles HTTP requests to get a specific stock
func GetStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockId := vars["id"]

	id, err := strconv.Atoi(stockId)
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	stock, err := db.DbGetStock(int64(id))
	if err != nil {
		http.Error(w, "Error fetching stock", http.StatusInternalServerError)
		return
	}

	if stock.StockID == 0 {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := db.DbGetAllStocks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stocks)
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	insertId, err := db.DbInsertStock(stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := Response{
		ID:      insertId,
		Message: "Stock created successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	stockId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedRows, err := db.DbUpdateStock(int64(stockId), stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	msg := fmt.Sprintf("Stock updated successfully. Rows affected: %d", updatedRows)

	res := Response{
		ID:      int64(stockId),
		Message: msg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockId, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	affectedRows, err := db.DbDeleteStock(int64(stockId))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	msg := fmt.Sprintf("Stock deleted successfully. Rows affected: %d", affectedRows)

	res := Response{
		ID:      int64(stockId),
		Message: msg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
