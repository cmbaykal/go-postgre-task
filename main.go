package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cmbaykal/go-postgre-task/database"
	"github.com/cmbaykal/go-postgre-task/models"
	"github.com/cmbaykal/go-postgre-task/routes"
	"github.com/gorilla/mux"
)

func main() {
	database.Connect()
	database.Db.AutoMigrate(&models.Ticket{})

	r := mux.NewRouter()
	r.HandleFunc("/ticket_options", routes.CreateTicket).Methods("POST")
	r.HandleFunc("/ticket/{id}", routes.GetTicket).Methods("GET")
	r.HandleFunc("/ticket_options/{id}/purchase", routes.PurchaseTicket).Methods("POST")

	fmt.Println("Starting Server port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
