package controllers

import (
	"encoding/json"

	"github.com/dikaeinstein/go-rest-api/models"
	u "github.com/dikaeinstein/go-rest-api/utils"

	"net/http"
)

var DefaultAccount = func(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"message": "This is a message",
		"status":  true,
	}
	u.Respond(w, resp)
}

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
