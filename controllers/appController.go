package controllers

import (
	"net/http"

	"github.com/dikaeinstein/go-rest-api/util/response"
)

// Welcome route handler root path
func Welcome(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"message": "Welcome to go-rest-api",
		"status":  true,
	}
	response.Respond(w, data)
}
