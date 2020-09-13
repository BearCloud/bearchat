package api

import (
	"fmt",
	"log",
	"net/http",
	"encoding/json",
	"errors",
	"database/sql",
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"time"
)

func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/posts", getFeed).Methods(http.MethodGet)
	router.HandleFunc("/api/posts/{uuid}", getPosts).Methods(http.MethodGet)
	router.HandleFunc("/api/posts", createPost).Methods(http.MethodPost)
	router.HandleFunc("/api/posts", deletePost).Methods(http.MethodDelete)

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
		postID string
		uuid string
		postTime Time
	)
	postsArray := make([]Post, counter)
	for i := 0; i < counter; i++ {
		err = posts.Scan(&content, &postID, &uuid, &postTime)
		if err != null {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
		postsArray[i] := Post{content, postID, uuid, postTime}
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
	var post string
	json.NewEncoder(w).Decode(&post)
	uuid, err := GetUUID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	postID := uuid.New().String()
	_, err := DB.Query("INSERT INTO posts(content, postID, uuid, postTime) VALUES (?, ?, ?, ?)", post, postID, uuid, time.Now())
	if err != nil {
		http.Error(w, errors.New("error storing post into database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	var postID string
	json.NewEncoder(w).Decode(&postID)
	uuid, err := GetUUID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	var exists bool
	//check if the email or username exists
	err = DB.QueryRow("SELECT EXISTS (SELECT * FROM posts WHERE postID = ?)", postID).Scan(&exists)
	if err != nil {
		http.Error(w, errors.New("error checking if post exists").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	if !exists {
		http.Error(w, errors.New("this post doesn't exist").Error(), http.StatusNotFound)
		return
	}
	var postUUID string
	_, err = DB.QueryRow("SELECT uuid FROM posts WHERE postID = ?", postID).Scan(&postUUID)
	if err != nil {
		http.Error(w, errors.New("error fetching post to delete from database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
	if uuid != postUUID {
		http.Error(w, errors.New("You are not authorized to delete this post").Error(), http.StatusUnauthorized)
		return
	}
	err = DB.QueryRow("DELETE FROM posts WHERE postID = ?", postID)
	if err != nil {
		http.Error(w, errors.New("error deleting post").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}
