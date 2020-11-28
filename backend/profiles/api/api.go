package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/profile/{uuid}", getProfile).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/profile/{uuid}", setProfile).Methods(http.MethodPut, http.MethodOptions)

	return nil
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	//check auth
	//fetch cookie
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
	}
	//validate the cookie
	claims, err := ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	log.Println(claims)
	//fetch public vs private depending on if user is accessing own profile
	var (
		first  string
		last   string
		email  string
		userid string
	)
	var profile Profile
	err = DB.QueryRow("SELECT * FROM users WHERE uuid = ?", uuid).Scan(&first, &last, &email, &userid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
	profile = Profile{first, last, email, userid}
	//to add later - more data if friends

	//encode fetched data as json and serve to client
	json.NewEncoder(w).Encode(profile)
}

func setProfile(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	//check auth - should also check if profile exists cause token is invalidated when profile deleted
	//fetch cookie
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
	}
	//validate the cookie
	claims, err := ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	// log.Println(claims)
	auth := (claims["UserID"] == uuid)

	if !auth {
		http.Error(w, errors.New("you are not authorized to edit this profile").Error(), http.StatusUnauthorized)
		log.Print(err.Error())
		return
	}
	//check for duplicate
	var created bool
	err = DB.QueryRow("SELECT EXISTS (SELECT UUID FROM users WHERE UUID = ?)", uuid).Scan(&created)
	//store new profile data if auth correct
	var profile Profile
	err = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	log.Println(profile)
	_, err = DB.Query("REPLACE INTO users(Firstname, Lastname, Email, UUID) VALUES (?, ?, ?, ?)", profile.Firstname, profile.Lastname, profile.Email, profile.UUID)
	if err != nil {
		http.Error(w, errors.New("error storing profile into database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}
