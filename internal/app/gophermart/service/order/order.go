package order

import (
	"context"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

type Service interface {
	Create(ctx context.Context, order model.Order) error
	UserOrders(ctx context.Context, user model.User) (model.Orders, error)
}

func NewService(repo repository.OrderRepository) Service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo repository.OrderRepository
}

func (s *service) Create(ctx context.Context, order model.Order) error {
	return nil
}

func (s *service) UserOrders(ctx context.Context, user model.User) (model.Orders, error) {
	return nil, nil
}
