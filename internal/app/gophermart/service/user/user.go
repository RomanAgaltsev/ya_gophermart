package user

import (
    "context"

    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var _ app.UserService = (*Service)(nil)

type Repository interface {
    CreateUser(ctx context.Context, user *model.User) error
    GetUser(ctx context.Context, login string) (*model.User, error)
}

func NewService(repository Repository) *Service {
    return &Service{
        repository: repository,
    }
}

type Service struct {
    repository Repository
}

func (s *Service) Register(ctx context.Context, user *model.User) error {
    return nil
}

func (s *Service) Login(ctx context.Context, user *model.User) error {
    return nil
}
