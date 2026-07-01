package auth

import (
	"github.com/HadeedTariq/go-production-grade-api/internal/utils/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauthConfig *oauth2.Config

func init() {
	// Initialize the OAuth2 configuration
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/github/callback",
		ClientID:     env.GetEnvString("GITHUB_CLIENT_ID", "client"),     // Set this in your terminal before running
		ClientSecret: env.GetEnvString("GITHUB_CLIENT_SECRET", "client"), // Set this in your terminal before running
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}
