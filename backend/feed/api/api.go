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
	router.HandleFunc("/api/feed/{uuid}", getFeed).Methods(http.MethodGet)

	return nil
}
