package main

import (
	"github.com/0-Inacio-0/API_TIS_IV/gyms"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
)


func main() {

	gyms.Init()

	router := gyms.NewRouter()
	// these two lines are important in order to allow access from the front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET"})


	log.Println(http.ListenAndServe(":8080", handlers.CORS(allowedOrigins, allowedMethods)(router)))

}
