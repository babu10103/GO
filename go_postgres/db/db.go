package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/babu10103/go/go-postgres/models"
	// "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {

	connStr := os.Getenv("POSTGRES_URL")
	if connStr == "" {
		return fmt.Errorf("POSTGRES_URL environment variable not set")
	}
	log.Println("Connection string:", connStr)

	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			// check if connection is functional and able to communicate with PostgreSQL server
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to the database!")
				return nil
			}
		}
		log.Printf("Attempt %d: Error connecting to PostgreSQL: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}
	return err
}

func GetDB() *sql.DB {
	return db
}

func DbGetStock(id int64) (models.Stock, error) {
	db := GetDB()

	var stock models.Stock
	query := "SELECT * FROM stocks WHERE stockid = $1"
	row := db.QueryRow(query, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	if err != nil {
		if err == sql.ErrNoRows {
			return stock, nil
		}
		return stock, fmt.Errorf("Unable to scan the row: %w", err)
	}

	return stock, nil
}

func DbGetAllStocks() ([]models.Stock, error) {
	db := GetDB()

	var stocks []models.Stock
	query := "SELECT * FROM stocks"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error querying database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating rows: %w", err)
	}

	return stocks, nil
}

func DbInsertStock(stock models.Stock) (int64, error) {
	db := GetDB()

	query := "INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid"

	var id int64

	err := db.QueryRow(query, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("Error inserting stock: %w", err)
	}
	return id, err
}

func DbDeleteStock(id int64) (int64, error) {
	db := GetDB()
	query := "DELETE FROM stocks WHERE stockid = $1"

	res, err := db.Exec(query, id)

	if err != nil {
		return -1, fmt.Errorf("Error deleting stock: %w", err)
	}

	rowsEffected, err := res.RowsAffected()

	if err != nil {
		return -1, fmt.Errorf("Error getting rows affected: %w", err)
	}

	return rowsEffected, nil
}

func DbUpdateStock(stockId int64, stock models.Stock) (int64, error) {
	db := GetDB()
	query := "UPDATE stocks SET name=$1, price=$2, company=$3 WHERE stockid=$4"

	res, err := db.Exec(query, stock.Name, stock.Price, stock.Company, stockId)

	if err != nil {
		return -1, fmt.Errorf("Error updating stock: %w", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return -1, fmt.Errorf("Error getting rows affected: %w", err)
	}
	return rowsAffected, nil
}
