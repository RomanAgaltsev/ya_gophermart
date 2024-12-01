package user

import (
    "context"

    "github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var _ Service = (*service)(nil)

type Repository interface {
    CreateUser(ctx context.Context, user *model.User) error
    GetUser(ctx context.Context, login string) (*model.User, error)
}

type Service interface {
    Register(ctx context.Context, user *model.User) error
    Login(ctx context.Context, user *model.User) error
}

func NewService(repository Repository) Service {
    return &service{
        repository: repository,
    }
}

type service struct {
    repository Repository
}

func (s *service) Register(ctx context.Context, user *model.User) error {
    return nil
}

func (s *service) Login(ctx context.Context, user *model.User) error {
    return nil
}
