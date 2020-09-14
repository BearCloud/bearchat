package api

import (
  "time"
)

type Post struct {
  Content string 'json:"content"'
  PostID string 'json:"postID"'
  Privacy bool 'json:"privacy"'
  UUID string 'json:"uuid"'
  PostTime Time 'json:"postTime"'
}
