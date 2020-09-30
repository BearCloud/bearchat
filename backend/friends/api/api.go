package api

import (
	"net/http"
	"github.com/go-gremlin/gremlin"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/friends/{uuid}", areFriends).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/friends/{uuid}", addFriend).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/friends/{uuid}", deleteFriend).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/friends/{uuid}/mutual", mutualFriends).Methods(http.MethodGet, http.MethodOptions)
	//router.HandleFunc("/api/friends", getFriends).Methods(http.MethodGet)
	router.HandleFunc("/api/friends", addUser).Methods(http.MethodPost, http.MethodOptions)

	return nil
}

func getUUID (w http.ResponseWriter, r *http.Request) (uuid string) {
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

	return claims["UserID"].(string)
}

func addUser (w http.ResponseWriter, r *http.Request) {
	uuid := getUUID(w, r)
	_, err := DB.Exec(gremlin.Query(`g.addV().property("uuid", userID)`).Bindings(gremlin.Bind{"userID": uuid}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}

func areFriends(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
	isFriend, err := DB.Exec(gremlin.Query(`g.V().has("uuid", userID).out("friends with").where(otherV().has("uuid", otherUUID)).count()`).Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	json.NewEncoder(w).Encode(string(isFriend))
}

func addFriend(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
	_, err := DB.Exec(gremlin.Query(`g.addEdge(g.V().has("uuid", userID).next(), g.V().has("uuid", otherUUID).next(), "friends with")`).Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

func deleteFriend(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
	_, err := DB.Exec(gremlin.Query(`g.V().has("uuid", userID).out("friends with").where(otherV().has("uuid", otherUUID)).drop()`).Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

func mutualFriends(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
	isFriend, err := DB.Exec(gremlin.Query(`g.V().has("uuid", userID).out("friends with").where(otherV().out("friends with").where(otherV().has("uuid", otherUUID)))`).Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	json.NewEncoder(w).Encode(isFriend)
}
