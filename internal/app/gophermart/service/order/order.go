package order

import (
	"context"

	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var _ Service = (*service)(nil)

type Repository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error)
}

type Service interface {
	Create(ctx context.Context, order *model.Order) error
	UserOrders(ctx context.Context, user *model.User) (model.Orders, error)
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

type service struct {
	repository Repository
}

func (s *service) Create(ctx context.Context, order *model.Order) error {
	return nil
}

func (s *service) UserOrders(ctx context.Context, user *model.User) (model.Orders, error) {
	return nil, nil
}
