package auth

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/HadeedTariq/go-production-grade-api/internal/utils/env"
	"github.com/gorilla/sessions"
)

var (
	SESSION_SECRET = []byte(env.GetEnvString("SESSION_SECRET", "hadeed@13"))
)

var Store = sessions.NewCookieStore(
	SESSION_SECRET,
)

const SessionName = "oauth-session"

func init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   false,
		SameSite: 2,
	}
}

func GenerateState() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
