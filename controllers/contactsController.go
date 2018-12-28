package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/dikaeinstein/go-rest-api/app"
	"github.com/dikaeinstein/go-rest-api/models"
	"github.com/dikaeinstein/go-rest-api/util/response"
)

// CreateContact creates a new user contact
func CreateContact(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) // Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		response.ErrorResponse(
			w,
			response.Message(false, "Error while decoding request body"),
			http.StatusBadRequest,
		)
		return
	}

	contact.UserID = user
	data := contact.Create()
	response.Respond(w, data)
}

// GetContactsFor retrieves all contacts for a specific user
func GetContactsFor(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(app.User("user")).(uint) // Grab the id of the user that sent the request

	contacts := models.GetContacts(user)
	data := response.Message(true, "success")
	data["data"] = contacts
	response.Respond(w, data)
}
