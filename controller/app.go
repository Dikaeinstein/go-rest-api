package controller

import (
	"net/http"

	"github.com/dikaeinstein/go-rest-api/util/response"
)

// Welcome route handler
func Welcome(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"message": "Welcome to go-rest-api",
		"status":  true,
	}
	response.Respond(w, data)
}

// NotFound route handler
func NotFound(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"message": "Not found",
		"status":  false,
	}
	response.Respond(w, data)
}
