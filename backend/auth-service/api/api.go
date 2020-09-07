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
	router.HandleFunc("/api/auth/signup", signup).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/signin", signin).Methods(http.MethodPost)
	router.HandleFunc("/api/auth/logout", logout).Methods(http.MethodPost)
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

	//obtain the credentials from the request body
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	//check if the email exists
	var exists bool
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
	_, err = DB.Query("INSERT INTO users(email, hashedPassword, verified, resetToken, userId, verifiedToken) VALUES (?, ?, FALSE, NULL, ?, ?)", credentials.Email, string(hashedPassword), userID, verificationToken)
	if err != nil {
		http.Error(w, errors.New("error storing credentials into database").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	// Set access token as a cookie
	var accessExpiresAt = time.Now().Add(DefaultAccessJWTExpiry)
	var accessToken string
	accessToken, err = setClaims(AuthClaims{
		Email:         credentials.Email,
		EmailVerified: false,
		UserID:        userID,
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
	})

	// Send verification email
	err = SendEmail(credentials.Email, "Email Verification", "user-signup.html", map[string]interface{}{"Token": verificationToken})
	if err != nil {
		http.Error(w, errors.New("error sending verification email").Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	return

}

func signin(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	var hashedPassword, userID string
	var verified bool
	err = DB.QueryRow("select hashedPassword, userId, verified from users where email=?", credentials.Email).Scan(&hashedPassword, &userID, &verified)
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
		Email:         credentials.Email,
		EmailVerified: verified,
		UserID:        userID,
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
	})

	return
}

func logout(w http.ResponseWriter, r *http.Request) {
	var expiresAt = time.Now().Add(-1 * time.Minute)
	http.SetCookie(w, &http.Cookie{Name: "access_token", Value: "", Expires: expiresAt})
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Expires: expiresAt})
	return
}
