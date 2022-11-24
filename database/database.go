package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Dsn = "host=localhost user=postgres password=1423337e dbname=postgres port=5432"
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
