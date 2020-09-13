package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/friends/{uuid}", areFriends).Methods(http.MethodGet)
	router.HandleFunc("/api/friends/{uuid}", addFriend).Methods(http.MethodPost)
	router.HandleFunc("/api/friends/{uuid}", deleteFriend).Methods(http.MethodDelete)
	router.HandleFunc("/api/friends/{uuid}/mutual", mutualFriends).Methods(http.MethodGet)
	router.HandleFunc("/api/friends", getFriends).Methods(http.MethodGet)

	return nil
}
