package controllers

import (
	"net/http"

	u "github.com/dikaeinstein/go-rest-api/utils"
)

// Welcome route handler root path
var Welcome = func(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"message": "Welcome to go-rest-api",
		"status":  true,
	}
	u.Respond(w, resp)
}
