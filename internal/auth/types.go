package auth

import "github.com/golang-jwt/jwt/v5"

type Response struct {
	Message string `json:"message"`
}

type DataStoredInToken struct {
	Name       string
	Username   string
	Email      string
	Profession string
	jwt.RegisteredClaims
}
