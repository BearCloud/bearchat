package api

import (
	"fmt",
	"log",
	"net/http",
	"encoding/json",
	"errors",
	"database/sql",
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/profile/{uuid}", fetchProfile).Methods(http.MethodGet)
	router.HandleFunc("/api/profile/{uuid}", setProfile).Methods(http.MethodPut)

	return nil
}

func checkAuth (r *http.Request, uuid string) (auth bool, err Error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	claim, err := getClaims(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	return (uuid == claim.UserID), err
}

func getProfile(w http.ResponseWriter, r *http.Request) {
  uuid := mux.Vars(r)["uuid"]
  //check auth
	auth, err := checkAuth(r, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
  //fetch public vs private depending on if user is accessing own profile
	var (
		first string
		last string
		email string
		dob string
		uuid string
	)
	var profile Profile
	if !auth {
		err := db.QueryRow("SELECT Firstname, Lastname FROM profiles WHERE uuid = ?", uuid).Scan(&first, &last)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
		profile = Profile{first, last, NULL, NULL, NULL}
	}
	else {
		err := db.QueryRow("SELECT * FROM profiles WHERE uuid = ?", uuid).Scan(&first, &last, &email, &dob, &uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
		profile = Profile{first, last, email, dob, uuid}
	}
	//to add later - more data if friends

  //encode fetched data as json and serve to client
  json.NewEncoder(w).Encode(profile)
}

func editProfile(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	//check auth - should also check if profile exists cause token is invalidated when profile deleted
	auth, err := checkAuth(r, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	if !auth {
		http.Error(w, errors.New("you are not authorized to edit this profile").Error(), http.StatusUnauthorized)
		log.Print(err.Error())
		return
	}
	//check for duplicate
	var created bool
	err = db.QueryRow("SELECT EXISTS (SELECT UUID FROM profiles WHERE UUID = ?)", uuid).Scan(&created)
	//store new profile data if auth correct
		profile := Profile{}
		err := json.NewDecoder(r.Body).Decode(&profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}
		if created {
			_, err = db.Query("REPLACE INTO profiles(Firstname, Lastname, Email, DOB, UUID) VALUES (?, ?, ?, ?, ?)", profile.Firstname, profile.Lastname, profile.Email, profile.DOB, profile.UUID)
		}
		else {
			_, err = db.Query("INSERT INTO profiles(Firstname, Lastname, Email, DOB, UUID) VALUES (?, ?, ?, ?, ?)", profile.Firstname, profile.Lastname, profile.Email, profile.DOB, profile.UUID)
		}
		if err != nil {
			http.Error(w, errors.New("error storing profile into database").Error(), http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}
}
