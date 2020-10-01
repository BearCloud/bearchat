package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"database/sql"
	"strconv"
)


func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/posts/{startIndex}", getFeed).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/posts/{uuid}/{startIndex}", getPosts).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/posts/create", createPost).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/posts/delete/{postID}", deletePost).Methods(http.MethodDelete, http.MethodOptions)

	return nil
}

func getUUID (w http.ResponseWriter, r *http.Request) (uuid string) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, errors.New("error obtaining cookie: " + err.Error()).Error(), http.StatusBadRequest)
		log.Print(err.Error())
		return ""
	}
	//validate the cookie
	claims, err := ValidateToken(cookie.Value)
	if err != nil {
		http.Error(w, errors.New("error validating token: " + err.Error()).Error(), http.StatusUnauthorized)
		log.Print(err.Error())
		return ""
	}
	log.Println(claims)

	return claims["UserID"].(string)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if (*r).Method == "OPTIONS" {
		return
	}

	uuid := mux.Vars(r)["uuid"]
	startIndex := mux.Vars(r)["startIndex"]
  	//check auth
	isAuthorized := (getUUID(w, r) == uuid)
	if !isAuthorized {
		http.Error(w, errors.New("Not authorized to get posts").Error(), http.StatusUnauthorized)
		return;
	}
  	//fetch public vs private depending on if user is accessing own profile
	var posts *sql.Rows
	var err error
	posts, err = DB.Query("SELECT * FROM posts WHERE authorID = ? ORDER BY postTime LIMIT ?, 25", uuid, startIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
<<<<<<< HEAD
<<<<<<< HEAD
	
	postID := uuid.New()
=======
=======
>>>>>>> 0f792e0fafe93ecd734de4d058d76046a9c4b1e6
	var (
		content string
		postID string
		userid string
		postTime time.Time
	)
	numPosts := 0
	postsArray := make([]Post, 25)
	for i := 0; i < 25 && posts.Next(); i++ {
		err = posts.Scan(&content, &postID, &userid, &postTime)
		if err != nil {
			http.Error(w, errors.New("Error scanning content: " + err.Error()).Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
		postsArray[i] = Post{content, postID, userid, postTime}
		numPosts++
	}
<<<<<<< HEAD
>>>>>>> master
=======
>>>>>>> 0f792e0fafe93ecd734de4d058d76046a9c4b1e6

	posts.Close()
	err = posts.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
  //encode fetched data as json and serve to client
  json.NewEncoder(w).Encode(postsArray[:numPosts])
  return;
}

func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if (*r).Method == "OPTIONS" {
		return
	}

	userID := getUUID(w, r)
	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	postID := uuid.New()
	pst, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = DB.Exec("INSERT INTO posts(content, postID, authorID, postTime) VALUES (?, ?, ?, ?)", post.PostBody, postID, userID, time.Now().In(pst))
	if err != nil {
		http.Error(w, errors.New("error storing post into database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if (*r).Method == "OPTIONS" {
		return
	}

	postID := mux.Vars(r)["postID"]
	//fetch cookie
	uuid := getUUID(w, r)
	log.Println(uuid)
	var exists bool
	//check if post exists
	err := DB.QueryRow("SELECT EXISTS (SELECT * FROM posts WHERE postID = ?)", postID).Scan(&exists)
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
	err = DB.QueryRow("SELECT authorID FROM posts WHERE postID = ?", postID).Scan(&postUUID)
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

func getFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if (*r).Method == "OPTIONS" {
		return
	}

	//get the start index
	startIndex := mux.Vars(r)["startIndex"]
	//convert to int
	intStartIndex, err := strconv.Atoi(startIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	//fetch cookie
	uuid := getUUID(w, r)
  //fetch public vs private depending on if user is accessing own profile
	posts, err := DB.Query("SELECT * FROM posts WHERE authorID <> ? ORDER BY postTime LIMIT ?, 25", uuid, intStartIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
	var (
		content string
		postID string
		userid string
		postTime time.Time
	)
	numPosts := 0
	postsArray := make([]Post, 25)
	for i := 0; i < 25 && posts.Next(); i++ {
		err = posts.Scan(&content, &postID, &userid, &postTime)
		if err != nil {
			http.Error(w, errors.New("Error scanning content: " + err.Error()).Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
		postsArray[i] = Post{content, postID, userid, postTime}
		numPosts++
	}

	err = posts.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
	}
  //encode fetched data as json and serve to client
  json.NewEncoder(w).Encode(postsArray[:numPosts])
}
