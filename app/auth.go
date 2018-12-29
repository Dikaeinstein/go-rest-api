package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dikaeinstein/go-rest-api/model"
	"github.com/dikaeinstein/go-rest-api/util/response"

	jwt "github.com/dgrijalva/jwt-go"
)

// User context key
type User string

// JwtAuthentication is the JWT middleware
func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login", "/api"} // List of endpoints that doesn't require auth
		reqPath := r.URL.Path                                           // current request path

		// check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == reqPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		data := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") // Grab the token from the header

		if tokenHeader == "" { // Token is missing, returns with error code 403 Unauthorized
			data = response.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response.Respond(w, data)
			return
		}

		splitted := strings.Split(tokenHeader, " ") // The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			data = response.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response.Respond(w, data)
			return
		}

		tokenPart := splitted[1] // Grab the token part, what we are truly interested in
		tk := &model.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { // Malformed token, returns with http code 403 as usual
			data = response.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response.Respond(w, data)
			return
		}

		if !token.Valid { // Token is invalid, maybe not signed on this server
			data = response.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response.Respond(w, data)
			return
		}

		// Everything went well, proceed with the request and
		// set the caller to the user retrieved from the parsed token
		// fmt.Sprintf("User %d", tk.UserID) // Useful for monitoring
		ctx := context.WithValue(r.Context(), User("user"), tk.UserID)
		r = r.WithContext(ctx)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
