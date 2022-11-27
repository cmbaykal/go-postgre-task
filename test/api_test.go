package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cmbaykal/go-postgre-task/main/database"
	"github.com/cmbaykal/go-postgre-task/main/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var fakeTicket = models.Ticket{
	ID:         1,
	Name:       "Ticket Name",
	Desc:       "Ticket Description",
	Allocation: 2,
}

func FakeCreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket
	json.NewDecoder(r.Body).Decode(&ticket)

	if ticket.Name == "" && ticket.Desc == "" && ticket.Allocation == 0{
		w.Write([]byte("Body Error"))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		fmt.Println("Test Success" + ticket.Name)
		createdTicket := database.Db.Create(&ticket)
		err := createdTicket.Error

		if err != nil {
			w.Write([]byte(err.Error()))
		}

		json.NewEncoder(w).Encode(&ticket)
		w.WriteHeader(http.StatusOK)
	}
}

func FakeGetTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var ticket models.Ticket
	dbResult := database.Db.Where("id = ?", params["id"]).Find(&ticket)

	if ticket.ID == 0 || errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
		w.Write([]byte("Ticket not found"))
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		json.NewEncoder(w).Encode(&ticket)
		w.WriteHeader(http.StatusOK)
	}
}

func FakePurchaseTicket(w http.ResponseWriter, r *http.Request) {
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

func TestCreateTicketSuccess(t *testing.T) {
	database.Connect("host=localhost user=baikal password=12345678 dbname=postgres_test port=52296")
	database.Db.Exec("DROP TABLE IF EXISTS tickets")
	database.Db.AutoMigrate(&models.Ticket{})

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

	FakeCreateTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	var ticket models.Ticket
	json.NewDecoder(res.Body).Decode(&ticket)

	if !reflect.DeepEqual(ticket, fakeTicket) {
		t.Errorf("FAILED: expected %v, got %v\n", fakeTicket, ticket)
	}
}

func TestCreateTicketFailed(t *testing.T) {
	body := bytes.NewBufferString(`
		{
			"testName": "Ticket Name",
			"testDec": "Ticket Description",
			"testAllocation": 2
		}
	`)

	mockRequest := httptest.NewRequest(http.MethodPost, "/ticket_options", body)
	mockRequest.Header.Set("Content-Type", "application/json")
	mockWriter := httptest.NewRecorder()

	FakeCreateTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	response := string(bodyBytes)

	if !reflect.DeepEqual(response, "Body Error") {
		t.Errorf("FAILED: expected %v, got %v\n", "Body Error", response)
	}
}
func TestGetTicketSuccess(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/ticket/1", nil)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	FakeGetTicket(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	var ticket models.Ticket
	json.NewDecoder(res.Body).Decode(&ticket)

	if !reflect.DeepEqual(ticket, fakeTicket) {
		t.Errorf("FAILED: expected %v, got %v\n", fakeTicket, ticket)
	}
}

func TestGetTicketFailed(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/ticket/2", nil)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	vars := map[string]string{
		"id": "2",
	}

	request = mux.SetURLVars(request, vars)

	FakeGetTicket(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	response := string(bodyBytes)

	if !reflect.DeepEqual(response, "Ticket not found") {
		t.Errorf("FAILED: expected %v, got %v\n", "Ticket not found", response)
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

	FakePurchaseTicket(mockWriter, mockRequest)

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

	FakePurchaseTicket(mockWriter, mockRequest)

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
	request := httptest.NewRequest(http.MethodGet, "/ticket/1", nil)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	FakeGetTicket(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	var ticket models.Ticket
	json.NewDecoder(res.Body).Decode(&ticket)

	if !reflect.DeepEqual(ticket.Allocation, 0) {
		t.Errorf("FAILED: expected %v, got %v\n", 0, ticket.Allocation)
	}

	database.Disconnect()
}
