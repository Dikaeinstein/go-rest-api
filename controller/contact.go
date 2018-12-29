package controller

import (
	"encoding/json"
	"net/http"

	"github.com/dikaeinstein/go-rest-api/middleware"
	"github.com/dikaeinstein/go-rest-api/model"
	"github.com/dikaeinstein/go-rest-api/util/response"
)

// CreateContact creates a new user contact
func CreateContact(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) // Grab the id of the user that send the request
	contact := &model.Contact{}

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
	user := r.Context().Value(middleware.User("user")).(uint) // Grab the id of the user that sent the request

	contacts := model.GetContacts(user)
	data := response.Message(true, "success")
	data["data"] = contacts
	response.Respond(w, data)
}
