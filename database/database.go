package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Dsn = "host=172.20.0.2 user=postgres password=postgres dbname=postgres port=5434"
var Db *gorm.DB
var err error

func Connect() {
	Db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: Dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DB Connection")
	}
}
