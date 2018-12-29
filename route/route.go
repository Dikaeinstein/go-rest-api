package route

import (
	"net/http"

	"github.com/dikaeinstein/go-rest-api/controller"
	"github.com/dikaeinstein/go-rest-api/middleware"
	"github.com/gorilla/mux"
)

// Router is a *mux.Router that should be passed to the HTTP listener
var Router *mux.Router

func init() {
	Router = routes() // Initialize Router
}

func routes() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	api.Path("").HandlerFunc(controller.Welcome).Methods(http.MethodGet)
	api.Path("/user/new").HandlerFunc(controller.CreateAccount).Methods(http.MethodPost)
	api.Path("/user/login").HandlerFunc(controller.Authenticate).Methods(http.MethodPost)
	api.Path("/contacts/new").HandlerFunc(controller.CreateContact).Methods(http.MethodPost)
	api.Path("/me/contacts").HandlerFunc(controller.GetContactsFor).Methods(http.MethodGet)
	api.Path("/").HandlerFunc(controller.NotFound)

	api.Use(middleware.JwtAuthentication) // Attach JWT auth middleware

	return router
}
