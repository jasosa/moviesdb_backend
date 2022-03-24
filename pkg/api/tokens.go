package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jasosa/backend/models"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

// tmp
var validUser = models.User{
	ID:       10,
	Email:    "me@here.com",
	Password: "$2a$12$MAYkOm4QWArsV5kczREsT.u7TZjyO7srH1lY6CI29WTSSe3DHJaMG",
}

type Credentials struct {
	Username string `"json:email"`
	Password string `"json:password"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		ErrorJSON(w, errors.New("Unauthorized"))
		return
	}

	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		ErrorJSON(w, errors.New("Unauthorized"))
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(""))
	if err != nil {
		ErrorJSON(w, errors.New("error signin"))
		return
	}

	WriteJSON(w, http.StatusOK, string(jwtBytes), "response")
}
