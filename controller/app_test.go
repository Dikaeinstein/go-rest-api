package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dikaeinstein/go-rest-api/model"
)

type ResponseData struct {
	Status  bool
	Message string
}

func cleanUp() {
	d := model.GetDB()
	fmt.Println("Cleaning up...")
	d.DropTable(&model.Account{}, &model.Contact{})
}

func TestWelcomeHandler(t *testing.T) {
	defer cleanUp()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(http.MethodGet, "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Welcome)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := ResponseData{
		Message: "Welcome to go-rest-api",
		Status:  true,
	}
	responseBody := &ResponseData{}
	json.NewDecoder(rr.Body).Decode(responseBody)
	if *responseBody != expected {
		t.Errorf("handler returned unexpected body: %v want %v",
			*responseBody, expected)
	}
}

func TestNotFoundHandler(t *testing.T) {
	defer cleanUp()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotFound)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: %v; want %v",
			status, http.StatusNotFound)
	}
	// Check the response body is what we expect.
	expected := ResponseData{
		Message: "Not found",
		Status:  false,
	}
	responseBody := &ResponseData{}
	json.NewDecoder(rr.Body).Decode(responseBody)
	if *responseBody != expected {
		t.Errorf("handler returned unexpected body: %v want %v",
			*responseBody, expected)
	}
}
