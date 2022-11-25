package test

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func DbURL() string {
	err := godotenv.Load("../.env")

	if err != nil {
		fmt.Println("Can not read enviroment config")
	}

	dbHost := os.Getenv("DB_TEST_HOST")
	dbPort := os.Getenv("DB_TEST_PORT")
	dbName := os.Getenv("DB_TEST_NAME")
	dbUser := os.Getenv("DB_TEST_USER")
	dbPassword := os.Getenv("DB_TEST_PASSWORD")

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
}

func Connect() {
	Db, err = gorm.Open(postgres.Open(DbURL()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DB Connection")
	}
}