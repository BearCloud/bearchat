package main

import (
	"log"
	"net/http"

	"github.com/BearCloud/fa20-project-dev/backend/posts/api"
	"github.com/gorilla/mux"
)

func main() {
	//init db
	DB := api.InitDB()
	defer DB.Close()

	//ping the database to make sure it's up
	err := DB.Ping()
	if err != nil {
		panic(err.Error())
	}
	// Create a new mux for routing api calls
	router := mux.NewRouter()

	err = api.RegisterRoutes(router)
	if err != nil {
		log.Fatal("Error registering API endpoints")
	}

	log.Println("listening...")
	log.Fatal(http.ListenAndServe(":80", router))
}
