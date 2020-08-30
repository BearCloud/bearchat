package api

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

//AuthClaims represents the claims in the access token
type AuthClaims struct {
	Email         string
	EmailVerified bool
	UserID        string
	jwt.StandardClaims
}

func getClaims(tokenString string) (claims AuthClaims, Error error) {
	claims = AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return AuthClaims{}, err
	}
	if !token.Valid {
		return AuthClaims{}, errors.New("The given token is not valid")
	}
	return claims, nil
}
