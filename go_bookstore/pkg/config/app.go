package config

import (
	"fmt"
	"os"

	"github.com/babu10103/GO/go_bookstore/pkg/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db      *gorm.DB
	dialect = "mysql"
)

func Connect() {
	DB_HOST := os.Getenv("DB_HOST")
	DB_USERNAME := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_USERNAME, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	var err error
	db, err = gorm.Open(dialect, connStr)
	if err != nil {
		log.ErrorLogger.Printf("Error while connecting to db: %v\n", err)
		log.ErrorLogger.Printf("connStr: %s\n", connStr)
		return
	}

	log.InfoLogger.Println("connected to db...")
}

func GetDB() *gorm.DB {
	return db
}
