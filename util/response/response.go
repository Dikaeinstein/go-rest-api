package response

import (
	"encoding/json"
	"net/http"
)

// Message creates a map of message and status
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond writes the json data to the HTTP response writer.
// It defaults to the HTTP OK status (200)
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// RespondWithStatus writes the json data and HTTP status code to the HTTP response writer
func RespondWithStatus(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
