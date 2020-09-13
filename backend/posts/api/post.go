package api

import (
  "time"
)

type Privacy int32

const(
  OnlyMe Privacy = iota
  Friends Privacy = iota
  Public Privacy = iota
)

type Post struct {
  Content string 'json:"content"'
  UUID string 'json:"uuid"'
  PrivacyLevel Privacy 'json:"privacyLevel"'
  PostTime Time 'json:"postTime"'
}
