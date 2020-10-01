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
	router.Use(CORS)

	err = api.RegisterRoutes(router)
	if err != nil {
		log.Fatal("Error registering API endpoints")
	}

	log.Println("listening...")
	log.Fatal(http.ListenAndServe(":80", router))
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> 0f792e0fafe93ecd734de4d058d76046a9c4b1e6
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "Content-Type")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
		return
	})
<<<<<<< HEAD
>>>>>>> master
}
=======
}
>>>>>>> 0f792e0fafe93ecd734de4d058d76046a9c4b1e6
