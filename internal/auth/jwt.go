package auth

import (
	"fmt"
	"time"

	"github.com/HadeedTariq/go-production-grade-api/internal/utils/env"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_SECRET               = []byte(env.GetEnvString("JWT_SECRET", "hadeed@13"))
	JWT_ACCESS_TOKEN_SECRET  = []byte(env.GetEnvString("JWT_REFRESH_TOKEN_SECRET", "hadeed@13"))
	JWT_REFRESH_TOKEN_SECRET = []byte(env.GetEnvString("JWT_REFRESH_TOKEN_SECRET", "hadeed@13"))
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

func GenerateAccessAndRefreshToken(data DataStoredInToken) (TokenResponse, error) {
	accessClaims := AccessTokenClaims{
		ID:         data.Id,
		Name:       data.Name,
		Username:   data.Username,
		Email:      data.Email,
		Profession: data.Profession,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
		},
	}

	refreshClaims := RefreshTokenClaims{
		ID: data.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * 24 * time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString(JWT_ACCESS_TOKEN_SECRET)
	if err != nil {
		return TokenResponse{}, err
	}

	refreshTokenString, err := refreshToken.SignedString(JWT_REFRESH_TOKEN_SECRET)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWT_ACCESS_TOKEN_SECRET, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
