package api

import (
	"github.com/go-gremlin/gremlin",
	"fmt",
	"log",
	"net/http",
	"encoding/json",
	"errors",
	"github.com/BearCloud/fa20-project-dev/api",
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/profile", fetchProfile).Methods(http.MethodPost)

	return nil
}

func fetchProfile()
