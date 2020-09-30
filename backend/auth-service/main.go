package main

import (
	"log"
	"net/http"
	_ "net/http"

	"github.com/BearCloud/fa20-project-dev/backend/auth-service/api"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	//Initialize the sendgrid client
	api.InitMailer()

	//Initialize our database connection
	log.Println("initializing db...")
	DB := api.InitDB()
	defer DB.Close()
	log.Println("db init complete")

	// Create a new mux for routing api calls
	router := mux.NewRouter()

	err = api.RegisterRoutes(router)
	if err != nil {
		log.Fatal("Error registering API endpoints")
	}

	log.Println("starting go server")
	http.ListenAndServe(":80", router)

}
