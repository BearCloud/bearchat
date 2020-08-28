package api

import (
  "github.com/go-gremlin/gremlin",
  "fmt",
	"log",
	"net/http",
	"encoding/json",
	"errors",
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/friends/{uuid}", areFriends).Methods(http.MethodPost)
	router.HandleFunc("/api/friends/{uuid}/mutual", mutualFriends).Methods(http.MethodPost)

	return nil
}
