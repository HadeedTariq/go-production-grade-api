package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	repo "github.com/HadeedTariq/go-production-grade-api/internal/adapters/postgresql/sqlc"
	authDto "github.com/HadeedTariq/go-production-grade-api/internal/auth/dto"
	"github.com/HadeedTariq/go-production-grade-api/internal/utils"
	"github.com/HadeedTariq/go-production-grade-api/internal/utils/env"
	"github.com/golang-jwt/jwt/v5"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ~ so according to me all the lower level stuff of the controller is managed with in the service
type Service interface {
	RegisterUser(ctx context.Context, data authDto.SignupRequest) (message string, err error)
}
type svc struct {
	repo *repo.Queries
	db   *pgxpool.Pool
}

func NewService(repo *repo.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) RegisterUser(ctx context.Context, data authDto.SignupRequest) (string, error) {
	_, err := s.repo.FindUserByEmail(ctx, data.Email)

	if err == nil {
		return "", errors.New("User Already exist")
	}

	tokenData := DataStoredInToken{
		Name:       data.Name,
		Username:   data.Username,
		Email:      data.Email,
		Profession: data.Profession,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := GenerateToken(tokenData)

	if err != nil {
		return "", err
	}

	encryptedToken, err := EncryptToken(token)

	if err != nil {
		return "", err
	}

	magicLink := fmt.Sprintf("%s/auth/register?token=%s", env.GetEnvString("SERVER_DOMAIN", "http://localhost:3000"), encryptedToken)

	// ~ so over there have to integrate the password hashing logic
	hashedPassword, err := utils.HashPassword(data.Password)

	if err != nil {
		return "", err
	}

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	_, err = qtx.CreateMagicLink(ctx, repo.CreateMagicLinkParams{
		Email: data.Email,
		Token: encryptedToken,
	})

	if err != nil {
		return "", err
	}
	_, err = qtx.CreateUser(ctx, repo.CreateUserParams{
		Name:     data.Name,
		Username: data.Username,
		Email:    data.Email,
		Profession: pgtype.Text{
			String: data.Profession,
			Valid:  true,
		},
		UserPassword: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
	})

	if err != nil {
		return "", err
	}

	// ~ so over there now have to integrate the verfication email sending functionality
	emailResponse := utils.SendVerificationEmail(data.Email, magicLink)

	if emailResponse.Success == false {
		return "", emailResponse.Error
	}

	tx.Commit(ctx)

	return "Verification email sent", nil
}
