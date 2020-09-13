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
	router.HandleFunc("/api/posts", getFeed).Methods(http.MethodGet)
	router.HandleFunc("/api/posts/{uuid}", getPosts).Methods(http.MethodGet)
	router.HandleFunc("/api/posts", createPost).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{postID}", deletePost).Methods(http.MethodDelete)

	return nil
}

func checkAuth (r *http.Request, uuid string) (auth bool, err Error) {
	userID, err := getUUID(r)
	return (uuid == userID), err
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
		posts, err := DB.Query("SELECT * FROM posts WHERE uuid = ? ORDER BY postTime", uuid)
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
	//to add later - more data if friends

  //encode fetched data as json and serve to client
  json.NewEncoder(w).Encode(postsArray)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	json.NewEncoder(w).Decode(&post)
	uuid, err := GetUUID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	_, err := DB.query("INSERT INTO posts(content, uuid, privacy, postTime) VALUES (?, ?, ?, ?)", post.Content, uuid, post.PrivacyLevel, post.PostTime)
	if err != nil {
		http.Error(w, errors.New("error storing post into database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}
