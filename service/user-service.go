package service

import (
	"context"
	"errors"
	repository "go-chatapp/repository/generated"
	"go-chatapp/utils"
)

type User struct {
	Name     string
	Email    string
	Password string
}

type UserService struct {
	queries *repository.Queries
}

func NewUserService(queries *repository.Queries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) RegisterAcc(ctx context.Context, req User) (User, error) {
	row, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return User{}, err
	}
	if row.Email != "" {
		return User{}, errors.New("user already exists")
	}

	hashedPass := utils.HashPassword(req.Password)

	user, err := s.queries.CreateUser(ctx, repository.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPass,
	})
	if err != nil {
		return User{}, err
	}

	return User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
