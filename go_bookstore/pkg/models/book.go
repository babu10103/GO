package models

import (
	"github.com/babu10103/GO/go_bookstore/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:""json:"name"`
	Authon      string `json:"authon"`
	Publication string `json: "publication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

// func (b *Book) CreateBook() *Book {
// 	db
// }
