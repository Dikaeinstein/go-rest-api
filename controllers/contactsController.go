package controllers

import (
	"encoding/json"

	"github.com/dikaeinstein/go-rest-api/app"
	"github.com/dikaeinstein/go-rest-api/models"

	u "github.com/dikaeinstein/go-rest-api/utils"

	"net/http"
)

// CreateContact creates a new user contact
func CreateContact(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint) // Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		u.ErrorResponse(
			w,
			u.Message(false, "Error while decoding request body"),
			http.StatusBadRequest,
		)
		return
	}

	contact.UserID = user
	resp := contact.Create()
	u.Respond(w, resp)
}

// GetContactsFor retrieves all contacts for a specific user
var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(app.User("user")).(uint) // Grab the id of the user that sent the request

	data := models.GetContacts(user)
	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
