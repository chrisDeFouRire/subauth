package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidAuthorization = errors.New("invalid authorization Bearer token header")
	ErrInvalidToken         = errors.New("invalid bearer JWT token")
	publicKey               *rsa.PublicKey
)

type Claims struct {
	jwt.RegisteredClaims
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}

func parseBearerToken(authorization string) (string, error) {
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", ErrInvalidAuthorization
	}
	token := strings.TrimPrefix(authorization, "Bearer ")
	return token, nil
}

func parseJWTToken(tokenString string) (*Claims, error) {
	c := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &c, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != "RS256" {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return &c, nil
}

func init() {
	spublicKey := os.Getenv("PUBLIC_KEY")
	if spublicKey == "" {
		panic("PUBLIC_KEY environment variable not set")
	}
	var err error
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte("-----BEGIN CERTIFICATE-----\n" + spublicKey + "\n-----END CERTIFICATE-----"))
	if err != nil {
		panic(err)
	}
	log.Println("Public key loaded")
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
		claims, err := parseJWTToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Println(claims)

		w.WriteHeader(http.StatusNoContent)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
	log.Println("Server started")
}
