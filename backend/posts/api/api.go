package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)


func RegisterRoutes(router *mux.Router) error {
	// router.HandleFunc("/api/posts/{startIndex}", getFeed).Methods(http.MethodGet)
	// router.HandleFunc("/api/posts/{uuid}/{startIndex}", getPosts).Methods(http.MethodGet)
	router.HandleFunc("/api/posts/create", createPost).Methods(http.MethodPost)
	router.HandleFunc("/api/posts/{postID}", deletePost).Methods(http.MethodDelete)

	return nil
}

// func getPosts(w http.ResponseWriter, r *http.Request) {

// 	uuid := mux.Vars(r)["uuid"]
// 	startIndex := mux.Vars(r)["startIndex"]
//   //check auth
// 	auth, err := checkAuth(r, uuid)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		log.Print(err.Error())
// 	}
//   //fetch public vs private depending on if user is accessing own profile
// 	var posts *sql.Rows
// 	if !auth {
// 		posts, err = DB.Query("SELECT * FROM posts WHERE uuid = ? AND privacy = 0 ORDER BY postTime OFFSET ? FETCH FIRST 25 ROWS ONLY", uuid, startIndex)
// 	} else {
// 		posts, err = DB.Query("SELECT * FROM posts WHERE uuid = ? ORDER BY postTime OFFSET ? FETCH FIRST 25 ROWS ONLY", uuid, startIndex)
// 	}
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		log.Print(err.Error())
// 	}
// 	defer posts.Close()
// 	postsPointer := posts
// 	counter := 0
// 	for postsPointer.Next() {
// 		counter++
// 	}
// 	var (
// 		content string
// 		postID string
// 		privacy bool
// 		uuid string
// 		postTime Time
// 	)
// 	postsArray := make([]Post, counter)
// 	for i := 0; i < counter; i++ {
// 		err = posts.Scan(&content, &postID, &privacy, &uuid, &postTime)
// 		if err != null {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			log.Print(err.Error())
// 		}
// 		postsArray[i] = Post{content, postID, uuid, postTime}
// 	}
// 	err = posts.Err()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		log.Print(err.Error())
// 	}
//   //encode fetched data as json and serve to client
//   json.NewEncoder(w).Encode(postsArray)
// }

func createPost(w http.ResponseWriter, r *http.Request) {
	//fetch cookie
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

	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	userID := claims["UserID"]

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Print(err.Error())
	}
	postID := uuid.New()

	_, err = DB.Exec("INSERT INTO posts(content, postID, authorID, postTime) VALUES (?, ?, ?, ?)", post.Content, postID, userID, time.Now())
	if err != nil {
		http.Error(w, errors.New("error storing post into database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["postID"]
	//fetch cookie
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

	uuid := claims["UserID"]

	var exists bool
	//check if post exists
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
	err = DB.QueryRow("SELECT uuid FROM posts WHERE postID = ?", postID).Scan(&postUUID)
	if err != nil {
		http.Error(w, errors.New("error fetching post to delete from database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
	if uuid != postUUID {
		http.Error(w, errors.New("You are not authorized to delete this post").Error(), http.StatusUnauthorized)
		return
	}
	_, err = DB.Exec("DELETE FROM posts WHERE postID = ?", postID)
	if err != nil {
		http.Error(w, errors.New("error deleting post").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

// func getFeed(w http.ResponseWriter, r *http.Request) {
// 	startIndex := mux.Vars(r)["startIndex"]
// 	uuid, err := GetUUID(r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		log.Print(err.Error())
// 	}
//   //fetch public vs private depending on if user is accessing own profile
// 	posts, err := DB.Query("SELECT * FROM posts WHERE uuid <> ? AND privacy = 0 ORDER BY postTime OFFSET ? FETCH FIRST 25 ROWS ONLY", uuid, startIndex)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		log.Print(err.Error())
// 	}
// 	defer posts.Close()
// 	postsPointer := posts
// 	counter := 0
// 	for postsPointer.Next() {
// 		counter++
// 	}
// 	var (
// 		content string
// 		postID string
// 		privacy bool
// 		uuid string
// 		postTime Time
// 	)
// 	postsArray := make([]Post, counter)
// 	for i := 0; i < counter; i++ {
// 		err = posts.Scan(&content, &postID, &privacy, &uuid, &postTime)
// 		if err != null {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			log.Print(err.Error())
// 		}
// 		postsArray[i] = Post{content, postID, uuid, postTime}
// 	}
// 	err = posts.Err()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		log.Print(err.Error())
// 	}
//   //encode fetched data as json and serve to client
//   json.NewEncoder(w).Encode(postsArray)
// }
