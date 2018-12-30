package model

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dikaeinstein/go-rest-api/util/response"
	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/bcrypt"
)

// Token is JWT claims struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// Account is a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email" sql:"unique_index"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token" sql:"-"`
}

// CheckDuplicateAccounts Checks for duplicate emails
func (account *Account) CheckDuplicateAccounts() (map[string]interface{}, bool) {
	// Email must be unique
	temp := &Account{}

	// Check for duplicate emails
	GetDB().Table("accounts").Where("email = ?", account.Email).First(temp)
	if temp.Email != "" {
		return response.Message(false,
			"Email address already in use by another user."), false
	}
	return response.Message(true, "Requirement passed"), true
}

// Validate incoming user details.
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return response.Message(false, "Email address is required"), false
	}
	if len(account.Password) < 6 {
		return response.Message(false, "Password is required"), false
	}

	return response.Message(true, "Requirement passed"), true
}

// Create user account
func (account *Account) Create() (map[string]interface{}, int, bool) {
	if data, ok := account.Validate(); !ok {
		return data, http.StatusBadRequest, false
	}

	if data, ok := account.CheckDuplicateAccounts(); !ok {
		return data, http.StatusConflict, false
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return response.Message(false, "Failed to create account, connection error."),
			http.StatusInternalServerError, false
	}

	// Create new JWT token for the newly registered account
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	account.Token = tokenString

	account.Password = "" // delete password

	data := response.Message(true, "Account has been created")
	data["account"] = account
	return data, http.StatusCreated, true
}

// Login user using email and password
func (account *Account) Login() (map[string]interface{}, int, bool) {
	if data, ok := account.Validate(); !ok {
		return data, http.StatusBadRequest, false
	}
	acc := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(acc).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return response.Message(false, "Email address not found"),
				http.StatusNotFound, false
		}
		return response.Message(false, "Connection error. Please retry"),
			http.StatusInternalServerError, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(account.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		// Password does not match!
		return response.Message(false, "Invalid login credentials. Please try again"),
			http.StatusUnauthorized, false
	}
	// Worked! Logged In
	account.Password = ""

	// Create JWT token
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))
	account.Token = tokenString // Store the token in the response

	data := response.Message(true, "Logged In")
	data["account"] = account
	return data, http.StatusOK, true
}

// GetUser retrieves user account using user id
func GetUser(userID uint) *Account {
	account := &Account{}
	userNotFound := GetDB().Table("accounts").Where("id = ?", userID).
		First(account).RecordNotFound()

	if userNotFound {
		return nil
	}

	account.Password = ""
	return account
}
