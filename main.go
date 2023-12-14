package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidAuthorization = errors.New("Invalid authorization header")
	publicKey               string
)

func parseBearerToken(authorization string) (string, error) {
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("Invalid authorization header")
	}
	token := strings.TrimPrefix(authorization, "Bearer ")
	return token, nil
}

func parseJWTToken(tokenString string, publicKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Return the key used for signing the token
		return []byte(publicKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return token, nil
}

func init() {
	publicKey = os.Getenv("PUBLIC_KEY")
	if publicKey == "" {
		panic("PUBLIC_KEY environment variable not set")
	}
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		token, err := parseBearerToken(authorization)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Parse the JWT token
		claims, err := parseJWTToken(token, publicKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Println(claims)

	})

	http.ListenAndServe(":8080", r)
}
