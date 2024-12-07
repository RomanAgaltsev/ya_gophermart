package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var (
	_ Service    = (*service)(nil)
	_ Repository = (*repository.Repository)(nil)

	ErrOrderUploadedByThisLogin    = fmt.Errorf("order number has already been uploaded by this user")
	ErrOrderUploadedByAnotherLogin = fmt.Errorf("order number has already been uploaded by another user")
)

type Service interface {
	Create(ctx context.Context, order *model.Order) error
	UserOrders(ctx context.Context, user *model.User) (model.Orders, error)
}

type Repository interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error)
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

func (s *service) Create(ctx context.Context, order *model.Order) error {
	existingOrder, err := s.repository.CreateOrder(ctx, order)
	if errors.Is(err, repository.ErrConflict) {
		if existingOrder.Login == order.Login {
			return ErrOrderUploadedByThisLogin
		}

		return ErrOrderUploadedByAnotherLogin
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *service) UserOrders(ctx context.Context, user *model.User) (model.Orders, error) {
	return s.repository.GetListOfOrders(ctx, user)
}
