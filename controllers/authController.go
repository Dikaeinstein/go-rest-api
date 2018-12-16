package controllers

import (
	"encoding/json"

	"github.com/dikaeinstein/go-rest-api/models"
	u "github.com/dikaeinstein/go-rest-api/utils"

	"net/http"
)

// CreateAccount creates a new user account using request payload
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) // Decode the request body into struct and failed if any error occur
	if err != nil {
		u.ErrorResponse(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp := account.Create() // Create account
	u.Respond(w, resp)
}

// Authenticate user using email and password
func Authenticate(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) // decode the request body into struct and failed if any error occur
	if err != nil {
		u.ErrorResponse(w, u.Message(false, "Invalid request"), http.StatusBadRequest)
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
