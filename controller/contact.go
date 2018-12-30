package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dikaeinstein/go-rest-api/middleware"
	"github.com/dikaeinstein/go-rest-api/model"
	"github.com/dikaeinstein/go-rest-api/util/response"
)

// CreateContact creates a new user contact
func CreateContact(w http.ResponseWriter, r *http.Request) {
	// Grab the id of the user that send the request
	user := r.Context().Value(middleware.User("user")).(uint)
	contact := &model.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		log.Println(err)
		response.RespondWithStatus(
			w,
			response.Message(false, "Error while decoding request body"),
			http.StatusBadRequest,
		)
		return
	}

	contact.UserID = user
	data, status, ok := contact.Create()
	if !ok {
		response.RespondWithStatus(w, data, status)
	}
	response.RespondWithStatus(w, data, status)
}

// GetContactsFor retrieves all contacts for a specific user
func GetContactsFor(w http.ResponseWriter, r *http.Request) {
	// Grab the id of the user that sent the request
	user := r.Context().Value(middleware.User("user")).(uint)

	contacts, status, ok := model.GetContacts(user)
	if !ok && status == http.StatusInternalServerError {
		data := response.Message(false, "Connection error. Please retry")
		response.RespondWithStatus(w, data, status)
		return
	}
	data := response.Message(true, "Success")
	data["data"] = contacts
	response.Respond(w, data)
}
