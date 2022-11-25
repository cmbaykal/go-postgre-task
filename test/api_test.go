package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cmbaykal/go-postgre-task/main/models"
	"github.com/gorilla/mux"
)

var mockTicket = models.Ticket{
	ID:         1,
	Name:       "Ticket Name",
	Desc:       "Ticket Description",
	Allocation: 2,
}

func MockCreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket
	json.NewDecoder(r.Body).Decode(&ticket)
	createdTicket := Db.Create(&ticket)
	err := createdTicket.Error

	if err != nil {
		w.Write([]byte(err.Error()))
	}

	json.NewEncoder(w).Encode(&ticket)
	w.WriteHeader(http.StatusOK)
}

func MockGetTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ticket models.Ticket
	Db.First(&ticket, params["id"])

	if ticket.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	json.NewEncoder(w).Encode(&ticket)
	w.WriteHeader(http.StatusOK)
}

func MockPurchaseTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ticket models.Ticket
	Db.First(&ticket, params["id"])

	var purchase models.TicketPurchase
	json.NewDecoder(r.Body).Decode(&purchase)

	if ticket.Allocation >= purchase.Quantity {
		Db.Model(&models.Ticket{}).Where("id = ?", ticket.ID).Update("allocation", ticket.Allocation-purchase.Quantity)

		w.Write([]byte("Purchase Complete"))
		w.WriteHeader(http.StatusOK)
	} else {
		w.Write([]byte("Not available ticket allocation"))
		w.WriteHeader(http.StatusNotFound)
	}
}

func TestCreateTicket(t *testing.T) {
	Connect()
	Db.Exec("DROP TABLE IF EXISTS tickets")
	Db.AutoMigrate(&models.Ticket{})

	body := bytes.NewBufferString(`
		{
			"name": "Ticket Name",
			"desc": "Ticket Description",
			"allocation": 2
		}
	`)

	mockRequest := httptest.NewRequest(http.MethodPost, "/ticket_options", body)
	mockRequest.Header.Set("Content-Type", "application/json")
	mockWriter := httptest.NewRecorder()

	MockCreateTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	var ticket models.Ticket
	json.NewDecoder(res.Body).Decode(&ticket)

	if !reflect.DeepEqual(ticket, mockTicket) {
		t.Errorf("FAILED: expected %v, got %v\n", mockTicket, ticket)
	}
}

func TestGetCreatedTicket(t *testing.T) {
	mockRequest := httptest.NewRequest(http.MethodGet, "/ticket/1", nil)
	mockRequest.Header.Set("Content-Type", "application/json")
	mockWriter := httptest.NewRecorder()

	MockGetTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	var ticket models.Ticket
	json.NewDecoder(res.Body).Decode(&ticket)

	if !reflect.DeepEqual(ticket, mockTicket) {
		t.Errorf("FAILED: expected %v, got %v\n", mockTicket, ticket)
	}
}

func TestTicketPurchaseSuccessful(t *testing.T) {
	body := bytes.NewBufferString(`
		{
			"quantity": 2,
			"user_id": "123456"
		}
	`)

	mockRequest := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchase", body)
	mockRequest.Header.Set("Content-Type", "application/json")
	mockWriter := httptest.NewRecorder()

	MockPurchaseTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

    if err != nil {
        t.Errorf(err.Error())
    }

	response := string(bodyBytes)

	if !reflect.DeepEqual(response, "Purchase Complete") {
		t.Errorf("FAILED: expected %v, got %v\n", "Purchase Complete", response)
	}
}

func TestTicketPurchaseFailed(t *testing.T) {
	body := bytes.NewBufferString(`
		{
			"quantity": 2,
			"user_id": "123456"
		}
	`)

	mockRequest := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchase", body)
	mockRequest.Header.Set("Content-Type", "application/json")
	mockWriter := httptest.NewRecorder()

	MockPurchaseTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

    if err != nil {
        t.Errorf(err.Error())
    }

	response := string(bodyBytes)

	if !reflect.DeepEqual(response, "Not available ticket allocation") {
		t.Errorf("FAILED: expected %v, got %v\n", "Not available ticket allocation", response)
	}
}

func TestTicketAllocation(t *testing.T) {
	mockRequest := httptest.NewRequest(http.MethodGet, "/ticket/1", nil)
	mockRequest.Header.Set("Content-Type", "application/json")
	mockWriter := httptest.NewRecorder()

	MockGetTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	var ticket models.Ticket
	json.NewDecoder(res.Body).Decode(&ticket)

	if !reflect.DeepEqual(ticket.Allocation, 0) {
		t.Errorf("FAILED: expected %v, got %v\n", 0, ticket.Allocation)
	}
}
