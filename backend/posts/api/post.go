package api

import "time"

type Post struct {
	Content  string    `json:"content"`
	PostID   string    `json:"postID"`
	AuthorID string    `json:"AuthorID"`
	PostTime time.Time `json:"postTime"`
}
