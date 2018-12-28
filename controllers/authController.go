package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dikaeinstein/go-rest-api/models"
	"github.com/dikaeinstein/go-rest-api/util/response"
)

// CreateAccount creates a new user account using request payload
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) // Decode the request body into struct and failed if any error occur
	if err != nil {
		response.ErrorResponse(w, response.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	data := account.Create() // Create account
	response.Respond(w, data)
}

// Authenticate user using email and password
func Authenticate(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) // decode the request body into struct and failed if any error occur
	if err != nil {
		response.ErrorResponse(w, response.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	data := models.Login(account.Email, account.Password)
	response.Respond(w, data)
}
