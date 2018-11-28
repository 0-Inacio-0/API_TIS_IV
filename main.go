package main

import (
	"github.com/0-Inacio-0/API_TIS_IV/gyms"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting HTTP Server... ")

	router := gyms.NewRouter()
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST"})

	//get the env port
	port := gyms.DetermineListenAddress()

	log.Println("Running on PORT", port)
	log.Println(http.ListenAndServe(port, handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
