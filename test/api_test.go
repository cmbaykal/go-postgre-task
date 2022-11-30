package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/cmbaykal/go-postgre-task/main/database"
	"github.com/cmbaykal/go-postgre-task/main/models"
	"github.com/cmbaykal/go-postgre-task/main/routes"
	"github.com/gorilla/mux"
)

var fakeTicket = models.Ticket{
	ID:         1,
	Name:       "Ticket Name",
	Desc:       "Ticket Description",
	Allocation: 2,
}

func TestCreateTicketSuccess(t *testing.T) {
	database.Connect("host=localhost user=baikal password=12345678 port=55871")
	database.Db.Exec("CREATE DATABASE postgres_test")
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

	routes.CreateTicket(mockWriter, mockRequest)

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

	routes.CreateTicket(mockWriter, mockRequest)

	res := mockWriter.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	response := string(bodyBytes)
	response = strings.TrimSuffix(response, "\n")

	if !reflect.DeepEqual(response, "Body Error") {
		t.Errorf("FAILED: expected %v, got %v\n", response, "Body Error")
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

	routes.GetTicket(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	response := string(bodyBytes)
	response = strings.TrimSuffix(response, "\n")

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

	request := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchase", body)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	routes.GetTicket(writer, request)

	res := writer.Result()
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

func TestTicketPurchaseFalseBody(t *testing.T) {
	body := bytes.NewBufferString(`
		{
			"testQuantity": 2,
			"testUserId": "123456"
		}
	`)

	request := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchase", body)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	routes.PurchaseTicket(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	response := string(bodyBytes)
	response = strings.TrimSuffix(response, "\n")

	if !reflect.DeepEqual(response, "Body Error") {
		t.Errorf("FAILED: expected %v, got %v\n", "Body Error", response)
	}
}

func TestTicketPurchaseFalseID(t *testing.T) {
	body := bytes.NewBufferString(`
		{
			"quantity": 2,
			"user_id": "123456"
		}
	`)

	request := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchase", body)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	vars := map[string]string{
		"id": "2",
	}

	request = mux.SetURLVars(request, vars)

	routes.PurchaseTicket(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf(err.Error())
	}

	response := string(bodyBytes)
	response = strings.TrimSuffix(response, "\n")

	if !reflect.DeepEqual(response, "Ticket not found") {
		t.Errorf("FAILED: expected %v, got %v\n", "Ticket not found", response)
	}
}

func TestTicketPurchaseFailed(t *testing.T) {
	body := bytes.NewBufferString(`
		{
			"quantity": 2,
			"user_id": "123456"
		}
	`)

	request := httptest.NewRequest(http.MethodPost, "/ticket_options/1/purchase", body)
	request.Header.Set("Content-Type", "application/json")
	writer := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	routes.PurchaseTicket(writer, request)

	res := writer.Result()
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

	routes.GetTicket(writer, request)

	res := writer.Result()
	defer res.Body.Close()

	var ticket models.Ticket
	json.NewDecoder(res.Body).Decode(&ticket)

	if !reflect.DeepEqual(ticket.Allocation, 0) {
		t.Errorf("FAILED: expected %v, got %v\n", 0, ticket.Allocation)
	}

	database.Disconnect()
}
