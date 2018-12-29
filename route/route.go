package route

import (
	"github.com/dikaeinstein/go-rest-api/app"
	"github.com/dikaeinstein/go-rest-api/controller"
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

	api.Path("").HandlerFunc(controller.Welcome).Methods("GET")
	api.Path("/user/new").HandlerFunc(controller.CreateAccount).Methods("POST")
	api.Path("/user/login").HandlerFunc(controller.Authenticate).Methods("POST")
	api.Path("/contacts/new").HandlerFunc(controller.CreateContact).Methods("POST")
	api.Path("/me/contacts").HandlerFunc(controller.GetContactsFor).Methods("GET")
	api.Path("/").HandlerFunc(controller.NotFound)

	api.Use(app.JwtAuthentication) // Attach JWT auth middleware

	return router
}
