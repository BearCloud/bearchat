package api

import (
  "time"
)

type Post struct {
  Content string 'json:"content"'
  PostID string 'json:"postID"'
  UUID string 'json:"uuid"'
  PostTime Time 'json:"postTime"'
}
