package api

import (
	"net/http"
	"github.com/go-gremlin/gremlin"
	"github.com/gorilla/mux"
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

func getUUID (r *http.Request) (uuid string, err Error) {
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
	return claim.UserID, err
}

func addUser (w http.ResponseWriter, r *http.Request) {
	uuid := getUUID(r)
	_, err := gremlin.Query('g.addV().property('uuid', userID)').Bindings(gremlin.Bind{"userID": uuid).Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}

func areFriends(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(r)
	isFriend, err := gremlin.Query('g.V().has('uuid', userID).out('friends with').where(otherV().has('uuid', otherUUID)).count()').Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}).Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	json.NewEncoder(w).Encode((isFriend != 0))
}

func addFriend(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(r)
	addFriend, err := gremlin.Query('g.addEdge(g.V().has('uuid', userID).next(), g.V().has('uuid', otherUUID).next(), 'friends with')').Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}).Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

func deleteFriend(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(r)
	delFriend, err := gremlin.Query('g.V().has('uuid', userID).out('friends with').where(otherV().has('uuid', otherUUID)).drop()').Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}).Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

func mutualFriends(w http.ResponseWriter, r *http.Request) {
	otherUUID := mux.Vars(r)["uuid"]
	uuid := getUUID(r)
	isFriend, err := gremlin.Query('g.V().has('uuid', userID).out('friends with').where(otherV().out('friends with').where(otherV().has('uuid', otherUUID)))').Bindings(gremlin.Bind{"userID": uuid, "otherUUID": otherUUID}).Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	json.NewEncoder(w).Encode(isFriend)
}
