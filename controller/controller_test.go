package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dikaeinstein/go-rest-api/middleware"
	"github.com/dikaeinstein/go-rest-api/model"
)

type ResponseData struct {
	Message string
	Status  bool
}

type handler func(w http.ResponseWriter, r *http.Request)

func createPayload(payload string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(payload))
}

func clearTable(table string) {
	d := model.GetDB()
	fmt.Printf("Cleaning up %s...", table)
	d.Exec(fmt.Sprintf("DELETE FROM %s", table))
	d.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table))
}

var appTests = []struct {
	route    string
	h        handler
	code     int
	expected ResponseData
}{
	{"/api", Welcome, http.StatusOK, ResponseData{"Welcome to go-rest-api", true}},
	{"/", NotFound, http.StatusNotFound, ResponseData{"Not found", false}},
}

func TestAppHandlers(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	for _, c := range appTests {
		req, err := http.NewRequest(http.MethodGet, c.route, nil)
		if err != nil {
			t.Fatal(err)
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.h)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != c.code {
			t.Errorf("handler returned wrong status code: %v; want %v",
				status, c.code)
		}

		// Check the response body is what we expect.
		responseBody := &ResponseData{}
		json.NewDecoder(rr.Body).Decode(responseBody)
		if *responseBody != c.expected {
			t.Errorf("handler returned unexpected body: %v; want %v",
				*responseBody, c.expected)
		}
	}
}

var authTests = []struct {
	route    string
	body     io.Reader
	h        handler
	code     int
	expected ResponseData
}{
	{"/api/user/new", createPayload(`{"email":"test","password":"password"}`),
		CreateAccount, http.StatusBadRequest,
		ResponseData{"Email address is required", false},
	},
	{"/api/user/new", createPayload(`{"email":"test@mail.com","password":""}`),
		CreateAccount, http.StatusBadRequest,
		ResponseData{"Password is required", false},
	},
	{"/api/user/new", createPayload(`{"email":"test@mail.com","password":"password"}`),
		CreateAccount, http.StatusCreated,
		ResponseData{"Account has been created", true},
	},
	{"/api/user/new", createPayload(`{"email":"test@mail.com","password":"password"}`),
		CreateAccount, http.StatusConflict,
		ResponseData{"Email address already in use by another user.", false},
	},
	{"/api/user/new", createPayload(`{"email":"test",password:"password"}`),
		CreateAccount, http.StatusBadRequest,
		ResponseData{"Invalid request", false},
	},
	{"/api/user/login", createPayload(`{"email":"test",password:"password"}`),
		Authenticate, http.StatusBadRequest,
		ResponseData{"Invalid request", false},
	},
	{"/api/user/login", createPayload(`{"email":"test","password":"password"}`),
		Authenticate, http.StatusBadRequest,
		ResponseData{"Email address is required", false},
	},
	{"/api/user/login", createPayload(`{"email":"test@mail.com","password":""}`),
		Authenticate, http.StatusBadRequest,
		ResponseData{"Password is required", false},
	},
	{"/api/user/login", createPayload(`{"email":"test123@mail.com","password":"password"}`),
		Authenticate, http.StatusNotFound,
		ResponseData{"Email address not found", false},
	},
	{"/api/user/login", createPayload(`{"email":"test@mail.com","password":"invalidpass"}`),
		Authenticate, http.StatusUnauthorized,
		ResponseData{"Invalid login credentials. Please try again", false},
	},
	{"/api/user/login", createPayload(`{"email":"test@mail.com","password":"password"}`),
		Authenticate, http.StatusOK, ResponseData{"Logged In", true},
	},
}

func TestAuthHandlers(t *testing.T) {
	clearTable("accounts")
	for _, c := range authTests {
		req, err := http.NewRequest(http.MethodPost, c.route, c.body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.h)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != c.code {
			t.Errorf("handler returned wrong status code: %v; want %v",
				status, c.code)
		}

		responseData := &ResponseData{}
		json.NewDecoder(rr.Body).Decode(responseData)
		if *responseData != c.expected {
			t.Errorf("handler returned unexpected body: %v; want %v",
				*responseData, c.expected)
		}
	}
}

var contactTests = []struct {
	method   string
	route    string
	body     io.Reader
	userID   uint
	h        handler
	code     int
	expected ResponseData
}{
	{http.MethodPost, "/api/contacts/new", createPayload(`{"phone":"08134568795"}`),
		1, CreateContact, http.StatusBadRequest,
		ResponseData{"Contact name should be on the payload", false},
	},
	{http.MethodPost, "/api/contacts/new", createPayload(`{"name":"Test User"}`),
		1, CreateContact, http.StatusBadRequest,
		ResponseData{"Phone number should be on the payload", false},
	},
	{http.MethodPost, "/api/contacts/new", createPayload(`{"name":"Test User",phone:"08134568795"}`),
		1, CreateContact, http.StatusBadRequest,
		ResponseData{"Error while decoding request body", false},
	},
	{http.MethodPost, "/api/contacts/new", createPayload(`{"name":"Test User","phone":"08134568795"}`),
		0, CreateContact, http.StatusBadRequest,
		ResponseData{"User is not recognized", false},
	},
	{http.MethodPost, "/api/contacts/new", createPayload(`{"name":"Test User","phone":"08134568795"}`),
		1, CreateContact, http.StatusCreated,
		ResponseData{"Successfully created contact", true},
	},
	{http.MethodGet, "/api/me/contacts", nil, 1, GetContactsFor, http.StatusOK,
		ResponseData{"Success", true},
	},
}

func TestContactHandlers(t *testing.T) {
	clearTable("contacts")
	for _, c := range contactTests {
		req, err := http.NewRequest(c.method, c.route, c.body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(c.h)

		// Populate the request's context with our test data.
		ctx := req.Context()
		ctx = context.WithValue(ctx, middleware.User("user"), uint(c.userID))
		// Add our context to the request: note that WithContext returns a copy of
		// the request, which we must assign.
		req = req.WithContext(ctx)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != c.code {
			t.Errorf("handler returned wrong status code: %v; want %v",
				status, c.code)
		}

		responseData := &ResponseData{}
		json.NewDecoder(rr.Body).Decode(responseData)
		if *responseData != c.expected {
			t.Errorf("handler return unexpected body: %v; want %v",
				*responseData, c.expected)
		}
	}
}
