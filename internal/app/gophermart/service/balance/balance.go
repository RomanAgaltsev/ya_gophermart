package balance

import (
	"context"
	"fmt"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var (
	_ Service    = (*service)(nil)
	_ Repository = (*repository.Repository)(nil)

	ErrNotEnoughBalance = fmt.Errorf("not enough balance for withdrawal")
)

type Service interface {
	UserBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	BalanceWithdraw(ctx context.Context, user *model.User, orderNumber string, sum float64) error
	UserWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
}

type Repository interface {
	GetBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	Withdraw(ctx context.Context, user *model.User, orderNumber string, sum float64) error
	GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
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

func (s *service) UserBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
	return nil, nil
}

func (s *service) BalanceWithdraw(ctx context.Context, user *model.User, orderNumber string, sum float64) error {
	return s.repository.Withdraw(ctx, user, orderNumber, sum)
}

func (s *service) UserWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
	return s.repository.GetListOfWithdrawals(ctx, user)
}
