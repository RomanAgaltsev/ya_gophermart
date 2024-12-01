package order

import (
	"context"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var _ app.OrderService = (*Service)(nil)

type Repository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error)
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

type Service struct {
	repository Repository
}

func (s *Service) Create(ctx context.Context, order *model.Order) error {
	return nil
}

func (s *Service) UserOrders(ctx context.Context, user *model.User) (model.Orders, error) {
	return nil, nil
}
