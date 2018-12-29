package model

import (
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
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token" sql:"-"`
}

// Validate incoming user details.
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return response.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return response.Message(false, "Password is required"), false
	}

	// Email must be unique
	temp := &Account{}

	// Check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return response.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return response.Message(false, "Email address already in use by another user."), false
	}

	return response.Message(false, "Requirement passed"), true
}

// Create user account
func (account *Account) Create() map[string]interface{} {
	if data, ok := account.Validate(); !ok {
		return data
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return response.Message(false, "Failed to create account, connection error.")
	}

	// Create new JWT token for the newly registered account
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" // delete password

	data := response.Message(true, "Account has been created")
	data["account"] = account
	return data
}

// Login user using email and password
func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Message(false, "Email address not found")
		}
		return response.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
		return response.Message(false, "Invalid login credentials. Please try again")
	}
	// Worked! Logged In
	account.Password = ""

	// Create JWT token
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString // Store the token in the response

	data := response.Message(true, "Logged In")
	data["account"] = account
	return data
}

// GetUser retrieves user account using user id
func GetUser(userID uint) *Account {
	account := &Account{}
	GetDB().Table("accounts").Where("id = ?", userID).First(account)
	if account.Email == "" { // User not found!
		return nil
	}

	account.Password = ""
	return account
}
