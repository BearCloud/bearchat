package api

import "time"

type Post struct {
	PostBody  string    `json:"postBody"`
	PostID   string    `json:"postID"`
	UUID     string    `json:"uuid"`
	PostTime time.Time `json:"postTime"`
	PostAuthor string `json:"postAuthor"`
}
