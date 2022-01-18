package dao

import (
	"log"
	"sms/service/src/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var database *gorm.DB

func InitDB() {
	db, err := gorm.Open(config.DB.Driver, config.GetConnArgs())
	if err != nil {
		log.Fatal(err)
	}
	db.SingularTable(true)
	db.DB().SetMaxOpenConns(40)
	db.DB().SetMaxIdleConns(2)
	database = db
}

func CloseDB() {
	if database != nil {
		database.Close()
	}
}
