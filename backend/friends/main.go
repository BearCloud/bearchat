package main

import (
	"log"
	_ "log"
	"net/http"
	_ "net/http"

	"github.com/BearCloud/fa20-project-dev/backend/friends/api"
	"github.com/gorilla/mux"
)

func main() {

	// Create a new mux for routing api calls
	router := mux.NewRouter()

	err := api.RegisterRoutes(router)
	if err != nil {
		log.Fatal("Error registering API endpoints")
	}

	http.ListenAndServe(":8080", router)
}
