package auth

import "github.com/golang-jwt/jwt/v5"

type Response struct {
	Message string `json:"message"`
}

type DataStoredInToken struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Profession string `json:"profession"`
	Avatar     string `json:"avatar"`
	jwt.RegisteredClaims
}

type AccessTokenClaims struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Profession string `json:"profession"`

	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	ID int64 `json:"id"`

	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
