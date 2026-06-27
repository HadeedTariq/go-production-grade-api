package auth

import (
	"fmt"

	"github.com/HadeedTariq/go-production-grade-api/internal/utils/env"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_SECRET               = env.GetEnvString("JWT_SECRET", "hadeed@13")
	JWT_ACCESS_TOKEN_SECRET  = env.GetEnvString("JWT_REFRESH_TOKEN_SECRET", "hadeed@13")
	JWT_REFRESH_TOKEN_SECRET = env.GetEnvString("JWT_REFRESH_TOKEN_SECRET", "hadeed@13")
)

func GenerateToken(data DataStoredInToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*DataStoredInToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &DataStoredInToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(*DataStoredInToken); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
