package auth

import "github.com/golang-jwt/jwt/v5"

type Response struct {
	Message string `json:"message"`
}

type DataStoredInToken struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Profession string `json:"profession"`
	jwt.RegisteredClaims
}
