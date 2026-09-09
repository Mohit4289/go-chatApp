package service

import (
	"context"
	repository "go-chatapp/repository/generated"
)

type ContactService struct {
	queries *repository.Queries
}

func NewContactService(queries *repository.Queries) *ContactService {
	return &ContactService{
		queries: queries,
	}
}

func (s *ContactService) GetAllUser(ctx context.Context) (UserList, error) {
	row, err := s.queries.GetAllUser(ctx)
	if err != nil {
		return UserList{}, err
	}

	return UserList{
		Name: row,
	}, nil
}
