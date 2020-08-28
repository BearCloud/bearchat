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
	router.HandleFunc("/api/posts/{uuid}", fetchPosts).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{uuid}/feed", fetchFeed).Methods(http.MethodPost)

	return nil
}
