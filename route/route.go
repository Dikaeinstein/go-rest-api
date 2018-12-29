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

	router.Path("/api").HandlerFunc(controller.Welcome).Methods("GET")
	router.Path("/api/user/new").HandlerFunc(controller.CreateAccount).Methods("POST")
	router.Path("/api/user/login").HandlerFunc(controller.Authenticate).Methods("POST")
	router.Path("/api/contacts/new").HandlerFunc(controller.CreateContact).Methods("POST")
	router.Path("/api/me/contacts").HandlerFunc(controller.GetContactsFor).Methods("GET")
	router.Path("/*").HandlerFunc(controller.NotFound).Methods()

	router.Use(app.JwtAuthentication) // Attach JWT auth middleware

	return router
}
