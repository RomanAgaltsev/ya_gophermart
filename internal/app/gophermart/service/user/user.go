package user

import (
    "context"

    "github.com/RomanAgaltsev/ya_gophermart/internal/config"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var _ Service = (*service)(nil)

type Service interface {
    Register(ctx context.Context, user *model.User) error
    Login(ctx context.Context, user *model.User) error
}

type Repository interface {
    CreateUser(ctx context.Context, user *model.User) error
    GetUser(ctx context.Context, login string) (*model.User, error)
}

func NewService(repository Repository, cfg *config.Config) (Service, error) {
    return &service{
        repository: repository,
        cfg:        cfg,
    }, nil
}

type service struct {
    repository Repository
    cfg        *config.Config
}

func (s *service) Register(ctx context.Context, user *model.User) error {
    return nil
}

func (s *service) Login(ctx context.Context, user *model.User) error {
    return nil
}
