package database

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

func Connect(Dsn string) error {
	Db, err = gorm.Open(postgres.Open(Dsn), &gorm.Config{})
	if err != nil {
		return errors.New("connection failed")
	} else {
		return nil
	}
}

func Disconnect(){
	db, _ := Db.DB()
	db.Close()
}