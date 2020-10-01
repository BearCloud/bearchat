package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtTokenSize    = 128
	verifyTokenSize = 6
	resetTokenSize  = 6
)

// RegisterRoutes initializes the api endpoints and maps the requests to specific functions
func RegisterRoutes(router *mux.Router) error {
	router.HandleFunc("/api/auth/signup", signup).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/auth/signin", signin).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/auth/logout", logout).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/auth/verify", verify).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/auth/sendreset", sendReset).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/auth/resetpw", resetPassword).Methods(http.MethodPost, http.MethodOptions)

	// Load sendgrid credentials
	err := godotenv.Load()
	if err != nil {
		return err
	}

	sendgridKey = os.Getenv("SENDGRID_KEY")
	sendgridClient = sendgrid.NewSendClient(sendgridKey)
	return nil
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if (*r).Method == "OPTIONS" {
		return
	}

	//obtain the credentials from the request body
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	var exists bool
	//check if the email or username exists
	err = DB.QueryRow("SELECT EXISTS (SELECT email FROM users WHERE username = ?)", credentials.Username).Scan(&exists)
	if err != nil {
		http.Error(w, errors.New("error checking if username exists").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	if exists == true {
		http.Error(w, errors.New("this username is taken").Error(), http.StatusConflict)
		return
	}

	err = DB.QueryRow("SELECT EXISTS (SELECT email FROM users WHERE email = ?)", credentials.Email).Scan(&exists)
	if err != nil {
		http.Error(w, errors.New("error checking if email exists").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	if exists == true {
		http.Error(w, errors.New("this email is already associated with an account").Error(), http.StatusConflict)
		return
	}

	//hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, errors.New("error hashing password").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	//Create a new user UUID and convert it to string
	userID := uuid.New().String()

	//Create new verification token
	verificationToken := GetRandomBase62(jwtTokenSize)

	//Store credentials in database
	_, err = DB.Query("INSERT INTO users(username, email, hashedPassword, verified, resetToken, userId, verifiedToken) VALUES (?, ?, ?, FALSE, NULL, ?, ?)",
		credentials.Username, credentials.Email, string(hashedPassword), userID, verificationToken)
	if err != nil {
		http.Error(w, errors.New("error storing credentials into database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	// Set access token as a cookie
	var accessExpiresAt = time.Now().Add(DefaultAccessJWTExpiry)
	var accessToken string
	accessToken, err = setClaims(AuthClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "access",
			ExpiresAt: accessExpiresAt.Unix(),
			Issuer:    defaultJWTIssuer,
			IssuedAt:  time.Now().Unix(),
		},
	})
	if err != nil {
		http.Error(w, errors.New("error creating accessToken").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   accessToken,
		Expires: accessExpiresAt,
		// Secure: true,
		// HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
		Path: "/",
	})

	// Set refresh token as a cookie.
	var refreshExpiresAt = time.Now().Add(DefaultAccessJWTExpiry)
	var refreshToken string
	refreshToken, err = setClaims(AuthClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "refresh",
			ExpiresAt: refreshExpiresAt.Unix(),
			Issuer:    defaultJWTIssuer,
			IssuedAt:  time.Now().Unix(),
		},
	})

	if err != nil {
		http.Error(w, errors.New("error creating refreshToken").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshToken,
		Expires: refreshExpiresAt,
		Path: "/",
	})

	// Send verification email
	err = SendEmail(credentials.Email, "Email Verification", "user-signup.html", map[string]interface{}{"Token": verificationToken})
	if err != nil {
		http.Error(w, errors.New("error sending verification email").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}


	w.WriteHeader(201)
	return

}

func signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if (*r).Method == "OPTIONS" {
		return
	}

	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	// if there is an error processing the credentials, print that error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	var hashedPassword, userID string
	var verified bool
	err = DB.QueryRow("select hashedPassword, userId, verified from users where username=?", credentials.Username).Scan(&hashedPassword, &userID, &verified)
	// process errors associated with emails
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, errors.New("this email is not associated with an account").Error(), http.StatusNotFound)
		} else {
			http.Error(w, errors.New("error retrieving information with this email").Error(), http.StatusInternalServerError)
			log.Print(err.Error())
		}
		return
	}

	// Check if hashed password matches the one corresponding to the email
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		http.Error(w, errors.New("the password you've entered is incorrect").Error(), http.StatusUnauthorized)
		return
	}

	// Set access token as a cookie.
	var accessExpiresAt = time.Now().Add(DefaultAccessJWTExpiry)
	var accessToken string
	accessToken, err = setClaims(AuthClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "access",
			ExpiresAt: accessExpiresAt.Unix(),
			Issuer:    defaultJWTIssuer,
			IssuedAt:  time.Now().Unix(),
		},
	})
	if err != nil {
		http.Error(w, errors.New("error creating accessToken").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   accessToken,
		Expires: accessExpiresAt,
		Path: "/",
	})

	// Set refresh token as a cookie.
	var refreshExpiresAt = time.Now().Add(DefaultAccessJWTExpiry)
	var refreshToken string
	refreshToken, err = setClaims(AuthClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "refresh",
			ExpiresAt: refreshExpiresAt.Unix(),
			Issuer:    defaultJWTIssuer,
			IssuedAt:  time.Now().Unix(),
		},
	})

	if err != nil {
		http.Error(w, errors.New("error creating refreshToken").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshToken,
		Expires: refreshExpiresAt,
		Path: "/",
	})
  w.WriteHeader(200)
}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if (*r).Method == "OPTIONS" {
		return
	}

	// logging out causes expiration time of cookie to be set to now
	var expiresAt = time.Now().Add(-1 * time.Minute)
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", Expires: expiresAt})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Expires: expiresAt})
	return
}

func verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if (*r).Method == "OPTIONS" {
		return
	}

	token, ok := r.URL.Query()["token"]
	// check that valid token exists
	if !ok || len(token[0]) < 1 {
		http.Error(w, errors.New("Url Param 'token' is missing").Error(), http.StatusInternalServerError)
		log.Print(errors.New("Url Param 'token' is missing").Error())
		return
	}

	_, err := DB.Exec("UPDATE users SET verified=1 WHERE verifiedToken = ?", token[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
	}

	return
}


func sendReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if (*r).Method == "OPTIONS" {
		return
	}

	//email from body
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	if credentials.Email == "" {
		http.Error(w, errors.New("Email field is empty").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	//generate reset token
	token := GetRandomBase62(resetTokenSize)

	//insert reset token into user database
	_, err = DB.Query("UPDATE users SET resetToken=? WHERE email = ?", token, credentials.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
	}

	// Send verification email
	err = SendEmail(credentials.Email, "BearChat Password Reset", "password-reset.html", map[string]interface{}{"Token": token})
	if err != nil {
		http.Error(w, errors.New("error sending verification email").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	return
}

func resetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if (*r).Method == "OPTIONS" {
		return
	}
	
	//get token from query params
	token := r.URL.Query().Get("token")

	//get the username, email, and password from the body
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	if credentials.Email == "" || credentials.Password == "" || credentials.Username == "" {
		http.Error(w, errors.New("Incorrect Credential Format").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	email := credentials.Email;
	username := credentials.Username;
	password := credentials.Password


	var exists bool
	//check if the token exists under the specified username
	err = DB.QueryRow("SELECT EXISTS (SELECT email FROM users WHERE username = ? AND resetToken=?)", username, token).Scan(&exists)
	if err != nil {
		http.Error(w, errors.New("error checking if token exists").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
	if exists == false {
		http.Error(w, errors.New("token doesnt exist").Error(), http.StatusConflict)
		return
	}


	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, errors.New("error hashing password").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	//input new password and clear the reset token
	_, err = DB.Exec("UPDATE users SET hashedPassword=?, resetToken=\"\" WHERE username=? and email=?", hashedPassword, username, email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print(err.Error())
	}

	//put the user in the redis cache to invalidate all current sessions


	return
}

