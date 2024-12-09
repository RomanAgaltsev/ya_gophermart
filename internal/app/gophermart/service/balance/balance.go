package balance

import (
	"context"
	"fmt"
	"time"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

const ordersProcessingInterval = 20

var (
	_ Service    = (*service)(nil)
	_ Repository = (*repository.Repository)(nil)

	ErrNotEnoughBalance = fmt.Errorf("not enough balance for withdrawal")
)

type Service interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, user *model.User) (*model.Balance, error)
	Withdraw(ctx context.Context, user *model.User, orderNumber string, sum float64) error
	Withdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
}

type Repository interface {
	CreateBalance(ctx context.Context, user *model.User) error
	GetBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	WithdrawFromBalance(ctx context.Context, user *model.User, orderNumber string, sum float64) error
	GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
	GetListOfOrdersToProcess(ctx context.Context) ([]string, error)
}

func NewService(repository Repository, cfg *config.Config) (Service, error) {
	balanceService := &service{
		repository: repository,
		cfg:        cfg,
	}

	go balanceService.ordersProcessing()

	return balanceService, nil
}

type service struct {
	repository Repository
	cfg        *config.Config
}

func (s *service) Create(ctx context.Context, user *model.User) error {
	return s.repository.CreateBalance(ctx, user)
}

func (s *service) Get(ctx context.Context, user *model.User) (*model.Balance, error) {
	return s.repository.GetBalance(ctx, user)
}

func (s *service) Withdraw(ctx context.Context, user *model.User, orderNumber string, sum float64) error {
	return s.repository.WithdrawFromBalance(ctx, user, orderNumber, sum)
}

func (s *service) Withdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
	return s.repository.GetListOfWithdrawals(ctx, user)
}

func (s *service) ordersProcessing() {
	ticker := time.NewTicker(ordersProcessingInterval * time.Second)

	for {
		select {
		case <-ticker.C:
			s.processOrders()
		default:
			continue
		}
	}
}

func (s *service) processOrders() {
	const workersNumber = 3

	//ordersToProcess := s.repository.

}
