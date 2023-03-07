package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtTokenExpiry = time.Minute * 15
const refreshTokenExpiry = time.Hour * 24

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

func (app *application) getTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// the authorization header looks like this:
	// Bearer <token>

	//add a header
	w.Header().Add("Vary", "Authorization")

	//get the authorization hehader
	authHeader := r.Header.Get("Authorization")

	//sanity check
	if authHeader == "" {
		return "", nil, errors.New("there is no auth header")
	}

	//split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	// check to see if we have the word "Bearer"
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("Unauthorized: no Bearer")
	}

	token := headerParts[1]

	// declare an empty Claims variable
	claims := &Claims{}

	// parse the token with the Claims using out secret
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		//validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("enexpected signing method: %v", token.Header["alg"])
		}

		return []byte(app.JWTSecret), nil
	})
	//check for an error; note that this catches expired tokens as well
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}

	//make sure that we issued this token
	if claims.Issuer != app.Domain {
		return "", nil, errors.New("incorrect issuer")
	}

	//token is valid
	return token, claims, nil
}
