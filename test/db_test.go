package test

import (
	"testing"

	"github.com/cmbaykal/go-postgre-task/main/database"
)

func TestDBConnectionSuccess(t *testing.T) {
	database.Connect("host=localhost user=baikal password=12345678 dbname=postgres_test port=52296")
}

func TestDBConnectionFailed(t *testing.T) {
	err := database.Connect("host=localhost user=baikal password=failed dbname=postgres_test port=52296")

	if err == nil{
		t.Errorf("FAILED: " + err.Error())
	}	
}