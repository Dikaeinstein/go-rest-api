package models

import (
	"fmt"

	u "github.com/dikaeinstein/go-rest-api/utils"

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
		return u.Message(false, "Contact name should be on the payload"), false
	}

	if contact.Phone == "" {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	if contact.UserID <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Create new contact
func (contact *Contact) Create() map[string]interface{} {

	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := u.Message(true, "success")
	resp["contact"] = contact
	return resp
}

// GetContact retrieve contact using contact id
func GetContact(id uint) *Contact {

	contact := &Contact{}
	err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error
	if err != nil {
		return nil
	}
	return contact
}

// GetContacts retrieves all the contacts that belongs to user
func GetContacts(user uint) []*Contact {

	contacts := make([]*Contact, 0)
	err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return contacts
}
