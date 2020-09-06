package api

import (
  "time"
)

type Profile struct {
  Firstname string 'json:"Firstname"'
  Lastname string 'json:"Lastname"'
  Email string 'json:"Email"'
  DOB Time 'json:"DOB"'
  UUID string 'json:"UUID"'
}
