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
	router.HandleFunc("/api/posts/{uuid}", getPosts).Methods(http.MethodGet)

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

func getPosts(w http.ResponseWriter, r *http.Request) {
  uuid := mux.Vars(r)["uuid"]
  //check auth
	auth, err := checkAuth(r, uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
  //fetch public vs private depending on if user is accessing own profile
	if !auth {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	else {
		posts, err := db.Query("SELECT * FROM posts WHERE uuid = ? ORDER BY postTime", uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
		defer posts.Close()
		postsPointer := posts
		counter := 0
		for postsPointer.Next() {
			counter++
		}
		var (
			content string
			uuid string
			privacylevel Privacy
			postTime Time
		)
		postsArray := make([]Post, counter)
		for i := 0; i < counter; i++ {
			err = posts.Scan(&content, &uuid, &privacylevel, &postTime)
			if err != null {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Print(err.Error())
			}
			postsArray[i] := Post{content, uuid, privacylevel, postTime}
		}
		err = posts.Err()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
	}
	//to add later - more data if friends

  //encode fetched data as json and serve to client
  json.NewEncoder(w).Encode(postsArray)
}
