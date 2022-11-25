package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cmbaykal/go-postgre-task/main/database"
	"github.com/cmbaykal/go-postgre-task/main/models"
	"github.com/cmbaykal/go-postgre-task/main/routes"
	"github.com/gorilla/mux"
)

func main() {
	database.Connect()
	database.Db.AutoMigrate(&models.Ticket{})

	r := mux.NewRouter()
	r.HandleFunc("/ticket_options", routes.CreateTicket).Methods(http.MethodPost)
	r.HandleFunc("/ticket/{id}", routes.GetTicket).Methods(http.MethodGet)
	r.HandleFunc("/ticket_options/{id}/purchase", routes.PurchaseTicket).Methods(http.MethodPost)

	fs := http.FileServer(http.Dir("../swaggerui"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	fmt.Println("Starting Api Server port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
