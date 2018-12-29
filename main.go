package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dikaeinstein/go-rest-api/model"
	"github.com/dikaeinstein/go-rest-api/route"
)

func main() {
	d := model.GetDB()
	defer d.Close() // Close db connection

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	fmt.Println("localhost:" + port)

	err := http.ListenAndServe(":"+port, route.Router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
