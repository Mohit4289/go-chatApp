package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"go-chatapp/config"
	repository "go-chatapp/repository/generated"
	"go-chatapp/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Name     string
	Email    string
	Password string
}

type UserResponse struct {
	ID    int64
	Name  string
	Email string
}

type UserService struct {
	queries *repository.Queries
}

func NewUserService(queries *repository.Queries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) RegisterAcc(ctx context.Context, req User) (UserResponse, string, string, error) {

	_, err := s.queries.GetUserByEmail(ctx, req.Email)

	if err == nil {
		return UserResponse{}, "", "", errors.New("user already exists")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return UserResponse{}, "", "", err
	}

	hashedPass := utils.HashPassword(req.Password)

	user, err := s.queries.CreateUser(ctx, repository.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPass,
	})
	if err != nil {
		return UserResponse{}, "", "", err
	}

	secretToken := config.EnvConfig().JWT_SECRET

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return UserResponse{}, "", "", errors.New("failed to generate access token")
	}

	generateToken := make([]byte, 32)
	if _, err := rand.Read(generateToken); err != nil {
		return UserResponse{}, "", "", errors.New("failed to generate refresh token")
	}

	stringToken := base64.RawURLEncoding.EncodeToString(generateToken)
	hash := sha256.Sum256([]byte(stringToken))
	tokenHash := hex.EncodeToString(hash[:])

	_, err = s.queries.AddRefreshToken(
		ctx,
		repository.AddRefreshTokenParams{
			RefreshToken: pgtype.Text{
				String: tokenHash,
				Valid:  true,
			},
			Email: user.Email,
		},
	)
	if err != nil {
		return UserResponse{}, "", "", err
	}

	return UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}, signedToken, stringToken, nil
}
