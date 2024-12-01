package balance

import (
	"context"

	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var _ Service = (*service)(nil)

type Repository interface {
	GetBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	Withdraw(ctx context.Context, user *model.User, order *model.Order, sum float64) error
	GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
}

type Service interface {
	UserBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	BalanceWithdraw(ctx context.Context, user *model.User, order *model.Order, sum float64) error
	UserWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

type service struct {
	repository Repository
}

func (s *service) UserBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
	return nil, nil
}

func (s *service) BalanceWithdraw(ctx context.Context, user *model.User, order *model.Order, sum float64) error {
	return nil
}

func (s *service) UserWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
	return nil, nil
}
