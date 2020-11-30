package api

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"encoding/json"
	"fmt"
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/friends/{uuid}", areFriends).Methods(http.MethodGet)
	router.HandleFunc("/api/friends/{uuid}", addFriend).Methods(http.MethodPost)
	router.HandleFunc("/api/friends/{uuid}", deleteFriend).Methods(http.MethodDelete)
	router.HandleFunc("/api/friends/{uuid}/mutual", mutualFriends).Methods(http.MethodGet)
	router.HandleFunc("/api/friends", getFriends).Methods(http.MethodGet)
	router.HandleFunc("/api/friends", addUser).Methods(http.MethodPost)

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

func getFriends (w http.ResponseWriter, r *http.Request) {
	uuid := getUUID(w, r)
	res, err := gremlinClient.Execute("g.V().has('uuid', '" + uuid + "').out('friends with').values('uuid')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}

	resJSON, err := json.Marshal(res[0].Result.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}

	fmt.Fprint(w, resJSON)
}

func addUser (w http.ResponseWriter, r *http.Request) {
	uuid := getUUID(w, r)
	_, err := gremlinClient.Execute("g.addV().property('uuid', '" + uuid + "')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}

func areFriends(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
  isFriend, err := gremlinClient.Execute("g.V().has('uuid', '" + uuid + "').out('friends with').where(otherV().has('uuid', '" + otherUUID + "')).count()")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	json.NewEncoder(w).Encode(string(isFriend[0].Result.Data))
}

func addFriend(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
	_, err := gremlinClient.Execute("g.addE('friends with').from(g.V().has('uuid', '" + uuid + "')).to(g.V().has('uuid', '" + otherUUID + "'))")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

func deleteFriend(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
  _, err := gremlinClient.Execute("g.V().bothE().filter(hasLabel('friends with')).where(inV().has('uuid', '" + uuid + "')).where(otherV().has('uuid', '" + otherUUID + "')).drop()")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

func mutualFriends(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(w, r)
	isFriend, err := gremlinClient.Execute("g.V().has('uuid', '" + uuid + "').both('friends with').and(both('friends with').has('uuid', '" +  otherUUID + "'))")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	json.NewEncoder(w).Encode(isFriend[0].Result.Data)
}
