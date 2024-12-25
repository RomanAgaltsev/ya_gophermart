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

// Service is the order service interface.
type Service interface {
	Create(ctx context.Context, order *model.Order) error
	UserOrders(ctx context.Context, user *model.User) (model.Orders, error)
}

// Repository is the order service repository interface.
type Repository interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error)
	GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error)
}

// NewService creates new order service.
func NewService(repository Repository, cfg *config.Config) (Service, error) {
	return &service{
		repository: repository,
		cfg:        cfg,
	}, nil
}

// service is the order service structure.
type service struct {
	repository Repository
	cfg        *config.Config
}

// Create creates new order.
func (s *service) Create(ctx context.Context, order *model.Order) error {
	// Create new order
	existingOrder, err := s.repository.CreateOrder(ctx, order)
	if errors.Is(err, repository.ErrConflict) {
		// There is a conflict - order number has been already uploaded
		if existingOrder.Login == order.Login {
			// By this user - user of request
			return ErrOrderUploadedByThisLogin
		}
		// By another user
		return ErrOrderUploadedByAnotherLogin
	}

	if err != nil {
		return err
	}

	return nil
}

// UserOrders returns a list of orders uploaded by user.
func (s *service) UserOrders(ctx context.Context, user *model.User) (model.Orders, error) {
	return s.repository.GetListOfOrders(ctx, user)
}
