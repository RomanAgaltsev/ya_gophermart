package user

import (
	"context"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

type Service interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) error
}

func NewService(repo repository.UserRepository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo repository.UserRepository
}

func (s *service) Register(ctx context.Context, user *model.User) error {
	return nil
}

func (s *service) Login(ctx context.Context, user *model.User) error {
	return nil
}
