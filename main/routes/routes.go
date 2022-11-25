package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cmbaykal/go-postgre-task/main/database"
	"github.com/cmbaykal/go-postgre-task/main/models"
	"github.com/gorilla/mux"
)

// swagger:operation POST /ticket_options postTicket
// ---
// produces:
// - application/json
// parameters:
//   - name: Body
//     in: body
//     description: Ticket options body for allocation
//     required: true
//     schema:
//     "$ref": "#/definitions/Ticket"
//
// responses:
//
//	'200':
//	  description: Created Ticket Body
//	  schema:
//	    "$ref": "#/definitions/Ticket"
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket
	json.NewDecoder(r.Body).Decode(&ticket)
	createdTicket := database.Db.Create(&ticket)
	err := createdTicket.Error

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(&ticket)
	w.WriteHeader(http.StatusOK)
}

// swagger:operation GET /ticket/{id} getTicket
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	  description: Found Ticket Body
//	  schema:
//	    "$ref": "#/definitions/Ticket"
func GetTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ticket models.Ticket
	database.Db.First(&ticket, params["id"])

	if ticket.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	json.NewEncoder(w).Encode(&ticket)
	w.WriteHeader(http.StatusOK)
}

// swagger:operation POST /ticket_options/{id}/purchase purchaseTicket
// ---
// produces:
// - application/json
// parameters:
//   - name: Body
//     in: body
//     description: Ticket Purchase body for purchase
//     required: true
//     schema:
//     "$ref": "#/definitions/Ticket"
//
// responses:
//
//	'200':
//	  description: Purchase Complete response
func PurchaseTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ticket models.Ticket
	database.Db.First(&ticket, params["id"])

	var purchase models.TicketPurchase
	json.NewDecoder(r.Body).Decode(&purchase)

	if ticket.Allocation >= purchase.Quantity {
		database.Db.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Update("allocation", ticket.Allocation-purchase.Quantity)

		w.Write([]byte("Purchase Complete"))
		w.WriteHeader(http.StatusOK)
	} else {
		w.Write([]byte("Not available ticket allocation"))
		w.WriteHeader(http.StatusNotFound)
	}
}
