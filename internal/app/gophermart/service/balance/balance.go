package balance

import (
	"context"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

var _ app.BalanceService = (*Service)(nil)

type Repository interface {
	GetBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	Withdraw(ctx context.Context, user *model.User, order *model.Order, sum float64) error
	GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

type Service struct {
	repository Repository
}

func (s *Service) UserBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
	return nil, nil
}

func (s *Service) BalanceWithdraw(ctx context.Context, user *model.User, order *model.Order, sum float64) error {
	return nil
}

func (s *Service) UserWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
	return nil, nil
}
