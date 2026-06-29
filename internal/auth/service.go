package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	middlewares "github.com/HadeedTariq/go-production-grade-api/internal"
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
	VerifyUser(ctx context.Context, token string) (message string, err error)
	LoginUser(ctx context.Context, data authDto.SigninRequest) (token TokenResponse, err error)
	AuthenticateUser(ctx context.Context) (user *AccessTokenClaims, err error)
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
		Id:         1,
		Avatar:     "",
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

	magicLink := fmt.Sprintf("%s/auth/verification?token=%s", env.GetEnvString("SERVER_DOMAIN", "http://localhost:3000"), encryptedToken)

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
	go utils.SendVerificationEmail(data.Email, magicLink)

	tx.Commit(ctx)

	return "Verification email sent", nil
}

func (s *svc) VerifyUser(ctx context.Context, token string) (string, error) {
	email, err := s.repo.FindMagicLinkByToken(ctx, token)
	if err != nil {
		return "", errors.New("invalid token")
	}

	decryptedToken, err := DecryptToken(token)
	if err != nil {
		return "", err
	}

	claims, err := ValidateToken(decryptedToken)
	if err != nil {
		return "", err
	}

	if claims.Email != email {
		return "", errors.New("incorrect token")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	user, err := qtx.VerifyUserByEmail(ctx, repo.VerifyUserByEmailParams{
		IsVerified: true,
		Email:      email,
	})

	if err != nil {
		return "", err
	}

	err = qtx.DeleteMagicLinksByEmail(ctx, claims.Email)
	if err != nil {
		return "", err
	}

	_, err = qtx.CreateAbout(ctx, user.ID)
	if err != nil {
		return "", err
	}

	_, err = qtx.CreateSocialLinks(ctx, user.ID)
	if err != nil {
		return "", err
	}

	_, err = qtx.CreateUserStats(ctx, user.ID)
	if err != nil {
		return "", err
	}

	_, err = qtx.CreateStreak(ctx, user.ID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return "User registered successfully", nil
}

func (s *svc) LoginUser(ctx context.Context, data authDto.SigninRequest) (token TokenResponse, err error) {
	user, err := s.repo.FindUserByEmail(ctx, data.Email)
	if err != nil {
		return TokenResponse{}, err
	}

	isPasswordCorrect := utils.CheckPasswordHash(data.Password, user.UserPassword.String)
	if !isPasswordCorrect {
		return TokenResponse{}, errors.New("incorrect credentials")
	}

	tokenResp, err := GenerateAccessAndRefreshToken(DataStoredInToken{
		Id:         user.ID,
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		Profession: user.Profession.String,
		Avatar:     user.Avatar.String,
	})
	if err != nil {
		return TokenResponse{}, err
	}

	err = s.repo.UpdateRefreshToken(ctx, repo.UpdateRefreshTokenParams{
		RefreshToken: pgtype.Text{
			String: tokenResp.RefreshToken,
			Valid:  true,
		},
		Email: user.Email,
	})

	if err != nil {
		return TokenResponse{}, err
	}

	return tokenResp, nil
}

func (s *svc) AuthenticateUser(ctx context.Context) (user *AccessTokenClaims, err error) {
	user = ctx.Value(middlewares.UserContextKey).(*AccessTokenClaims)
	return user, nil
}
