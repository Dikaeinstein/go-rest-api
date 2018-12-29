package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dikaeinstein/go-rest-api/route"
)

func main() {
	// Get port from .env file, we did not specify any port
	// so this should return an empty string when tested locally
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	fmt.Println("localhost:" + port)

	// Launch the app
	err := http.ListenAndServe(":"+port, route.Router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
