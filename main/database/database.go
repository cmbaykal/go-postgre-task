package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "host=postgres user=baikal password=12345678 dbname=postgres port=5434"
var Db *gorm.DB
var err error

func Connect() {
	Db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DB Connection")
	}
}
