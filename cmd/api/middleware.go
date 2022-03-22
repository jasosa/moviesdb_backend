package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		log.Println("Checkin token")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// could set anonymous user
		}

		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("invalid auth header"))
			return
		}

		if headerParts[0] != "Bearer" {
			app.errorJSON(w, errors.New("unauthorized - No Bearer"))
			return
		}

		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized - Failed HCMA Check"))
			return
		}

		err = validateClaims(claims)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized"))
			return
		}

		log.Println("Valid User:", userID)
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

func validateClaims(claims *jwt.Claims) error {

	if !claims.Valid(time.Now()) {
		return errors.New("unauthorized - expired token")
	}

	if !claims.AcceptAudience("mydomain.com") {
		return errors.New("unauthorized - invalid audience")
	}

	if claims.Issuer != "mydomain.com" {
		return errors.New("unauthorized - invalid issuer")
	}

	return nil
}
