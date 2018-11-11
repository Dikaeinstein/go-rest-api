package controllers

import (
	"encoding/json"

	"github.com/dikaeinstein/go-rest-api/models"

	u "github.com/dikaeinstein/go-rest-api/utils"

	"net/http"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	contact.UserId = user
	resp := contact.Create()
	u.Respond(w, resp)
}

var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) //Grab the id of the user that sent the request

	data := models.GetContacts(user)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
