package api

//Credentials respresents the user login object
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
