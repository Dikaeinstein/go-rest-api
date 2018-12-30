package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dikaeinstein/go-rest-api/model"
	"github.com/dikaeinstein/go-rest-api/util/response"
)

// CreateAccount creates a new user account using request payload
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := &model.Account{}
	// Decode the request body into struct and fail if any error occur
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		response.RespondWithStatus(w,
			response.Message(false, "Invalid request"),
			http.StatusBadRequest)
		return
	}

	data, status, ok := account.Create() // Create account
	if !ok {
		response.RespondWithStatus(w, data, status)
		return
	}
	response.RespondWithStatus(w, data, status)
}

// Authenticate user using email and password
func Authenticate(w http.ResponseWriter, r *http.Request) {
	account := &model.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		response.RespondWithStatus(w,
			response.Message(false, "Invalid request"),
			http.StatusBadRequest)
		return
	}

	data, status, ok := account.Login()
	if !ok {
		response.RespondWithStatus(w, data, status)
		return
	}
	response.Respond(w, data)
}
