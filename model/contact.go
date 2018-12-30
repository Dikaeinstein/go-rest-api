package model

import (
	"log"
	"net/http"

	"github.com/dikaeinstein/go-rest-api/util/response"
	"github.com/jinzhu/gorm"
)

// Contact model
type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserID uint   `json:"user_id"` // The user that this contact belongs to
}

// Validate method of Contact validates the required parameters
// in the request body. It returns the message and true if the requirements
// are met or false if otherwise.
func (contact *Contact) Validate() (map[string]interface{}, bool) {
	if contact.Name == "" {
		return response.Message(false, "Contact name should be on the payload"), false
	}

	if contact.Phone == "" {
		return response.Message(false, "Phone number should be on the payload"), false
	}

	if contact.UserID <= 0 {
		return response.Message(false, "User is not recognized"), false
	}
	// All the required parameters are present
	return response.Message(true, "Requirement passed"), true
}

// Create new contact
func (contact *Contact) Create() (map[string]interface{}, int, bool) {
	if data, ok := contact.Validate(); !ok {
		return data, http.StatusBadRequest, false
	}

	GetDB().Create(contact)

	data := response.Message(true, "Successfully created contact")
	data["contact"] = contact
	return data, http.StatusCreated, true
}

// GetContact retrieve contact using contact id
func GetContact(id uint) (*Contact, int, bool) {
	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, http.StatusNotFound, false
		}
		log.Println(err)
		return nil, http.StatusInternalServerError, false
	}
	return nil, http.StatusOK, true
}

// GetContacts retrieves all the contacts that belongs to user
func GetContacts(user uint) ([]*Contact, int, bool) {
	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Println(err)
		return nil, http.StatusInternalServerError, false
	}
	return contacts, http.StatusOK, true
}
