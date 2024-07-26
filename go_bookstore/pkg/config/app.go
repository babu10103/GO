package config

import (
	"fmt"
	// "net/url"

	"github.com/babu10103/GO/go_bookstore/pkg/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const DB_HOST = "10.125.56.126"
const DB_USERNAME = "root"
const DB_PASSWORD = "Babu@103"
const DB_NAME = "bookstore"
const DB_PORT = "3306"

var (
	db      *gorm.DB
	dialect = "mysql"
)

func Connect() {
	// encodedPassword := url.QueryEscape(DB_PASSWORD)
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_USERNAME, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	log.InfoLogger.Println(connStr)
	dbConnObj, err := gorm.Open(dialect, connStr)

	if err != nil {
		log.ErrorLogger.Printf("Error while connecting to db: %v\n", err)
		return
	}
	db = dbConnObj
}

func GetDB() *gorm.DB {
	return db
}
