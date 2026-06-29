package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/HadeedTariq/go-production-grade-api/internal/auth"
)

type ContextKey string

const UserContextKey ContextKey = "user"

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")

	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	cookie, err := r.Cookie("accessToken")
	if err == nil {
		return cookie.Value
	}

	return ""
}

// ~ ok so over there this is how the middleware is defined with in the golang
func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)

		if token == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		user, err := auth.ValidateAccessToken(token)

		if err != nil {
			http.Error(w, "Invalid or expired access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			UserContextKey,
			user,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
