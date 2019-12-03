package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"
)

func TestConn(t *testing.T) {
	_, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		t.Errorf("Connection Failed")
	}
}

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/result", GetResult).Methods("GET")
	router.HandleFunc("/api/result/{orderid}", GetSingleResult).Methods("GET")
	router.HandleFunc("/api/result", PostResult).Methods("POST")
	router.HandleFunc("/api/result/{orderid}", DeleteResult).Methods("DELETE")
	router.HandleFunc("/api/result/{orderid}", UpdateResult).Methods("PUT")

	return router

}
func TestGetResult(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/result", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Connection Failed, got:%v, want:%v", status, http.StatusOK)
	}
	dberror = true
	request1, _ := http.NewRequest("GET", "/api/result", nil)
	response1 := httptest.NewRecorder()
	Router().ServeHTTP(response1, request1)
	if status1 := response1.Code; status1 != http.StatusInternalServerError {
		t.Errorf("Want Error but got no Error!")
	}
	dberror = false
}

// Make a variable to a func and assign it to a negative of that function(which needs to created)
func TestGetVals(t *testing.T) {
	err := getVals()
	if err != nil {
		t.Errorf("error occur :%v", err)
	}
	findValFunc = findValN
	err = getVals()
	if err == nil {
		t.Error("No Error!")
	}
}

// How to create a function that creates an error
func findValN() error {
	return errors.New("error")
}

func TestGetSingleResult(t *testing.T) {
	request, _ := http.NewRequest("GET", "/api/result/1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Connection Failed, got:%v, want:%v", status, http.StatusOK)
	}
	request1, _ := http.NewRequest("GET", "/api/result/5", nil)
	response1 := httptest.NewRecorder()
	Router().ServeHTTP(response1, request1)
	if status := response1.Code; status != 404 {
		t.Error("Entry Exists")
	}
	request2, _ := http.NewRequest("GET", "/api/result/devansh", nil)
	response2 := httptest.NewRecorder()
	Router().ServeHTTP(response2, request2)
	if status := response2.Code; status != 500 {
		t.Error("Unknown Error!")
	}
}
func TestPostResult(t *testing.T) {
	p := Product{
		OrderID:  4,
		Name:     "asus",
		Price:    30000,
		Quantity: 1,
		Status:   true,
	}
	b, _ := json.Marshal(p)
	//	fmt.Println(bytes.NewBuffer(b))
	request, _ := http.NewRequest("POST", "/api/result", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Error("Post Unsuccessful")
	}
	request1, _ := http.NewRequest("POST", "/api/result", bytes.NewBuffer(b))
	response1 := httptest.NewRecorder()
	Router().ServeHTTP(response1, request1)
	if status := response1.Code; status != http.StatusConflict {
		t.Errorf("Unknown Error")
	}
	dberror = true
	request2, _ := http.NewRequest("POST", "/api/result", bytes.NewBuffer(nil))
	response2 := httptest.NewRecorder()
	Router().ServeHTTP(response2, request2)
	if status := response2.Code; status != 500 {
		t.Error("Unknown Error")
	}
	dberror = false
}

func TestUpdateResult(t *testing.T) {
	p := Product{
		OrderID:  4,
		Name:     "asus",
		Price:    30000,
		Quantity: 2,
		Status:   true,
	}
	b, _ := json.Marshal(p)
	request, _ := http.NewRequest("PUT", "/api/result/4", bytes.NewBuffer(b))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Connection Successful")
	} else if status == http.StatusConflict {
		t.Errorf("Entry already exists")
	}
}

func TestDeleteResult(t *testing.T) {
	request, _ := http.NewRequest("DELETE", "/api/result/4", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("Delete Succesful")
	} else if status == http.StatusNoContent {
		t.Errorf("Invalid Order ID!")
	}
}

func TestMain(t *testing.T) {
	go main()
}
