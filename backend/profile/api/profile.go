package api

type Profile struct {
  Firstname string 'json:"Firstname"'
  Lastname string 'json:"Lastname"'
  Email string 'json:"Email"'
  DOB string 'json:"DOB"'
  UUID string 'json:"UUID"'
}
