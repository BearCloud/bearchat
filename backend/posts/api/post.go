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
  Content string 'json:"Content"'
  UUID string 'json:"UUID"'
  PrivacyLevel Privacy 'json:"PrivacyLevel"'
  PostTime Time 'json:"PostTime"'
}
