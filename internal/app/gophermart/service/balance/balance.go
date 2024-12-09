package balance

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-chi/render"
)

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
	GetListOfOrdersToProcess(ctx context.Context) (model.Orders, error)
	UpdateBalanceAccrued(ctx context.Context, order *model.Order, accrual *model.OrderAccrual) error
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
	err := s.repository.WithdrawFromBalance(ctx, user, orderNumber, sum)
	if errors.Is(err, repository.ErrNegativeBalance) {
		return ErrNotEnoughBalance
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *service) Withdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
	return s.repository.GetListOfWithdrawals(ctx, user)
}

func (s *service) ordersProcessing() {
	const ordersProcessingInterval = 20

	ticker := time.NewTicker(ordersProcessingInterval * time.Second)

	for {
		select {
		case <-ticker.C:
			s.processOrders()
			// TODO
			//		default:
			//			continue
		}
	}
}

func (s *service) processOrders() {
	const workersNumber = 3

	ctx := context.Background()

	ordersToProcess, err := s.repository.GetListOfOrdersToProcess(ctx)
	if err != nil {
		slog.Info("orders processing", "error", err.Error())
		return
	}

	jobs := make(chan *model.Order, len(ordersToProcess))
	done := make(chan struct{}, len(ordersToProcess))

	for range workersNumber {
		go func(jobs chan *model.Order, done chan struct{}) {
			client := http.Client{}
			for order := range jobs {
				resp, errAccrual := backoff.RetryWithData(func() (*http.Response, error) {
					return client.Get(fmt.Sprintf("%s/api/orders/%s", s.cfg.AccrualSystemAddress, order.Number))
				}, backoff.NewExponentialBackOff())
				if errAccrual != nil {
					slog.Info("orders processing", "error", errAccrual.Error())
					done <- struct{}{}
					continue
				}

				var accrual model.OrderAccrual

				errRender := render.DecodeJSON(resp.Body, &accrual)
				if errRender != nil {
					slog.Info("orders processing", "error", errRender.Error())
					done <- struct{}{}
					continue
				}

				if order.Status == accrual.Status {
					done <- struct{}{}
					continue
				}

				errUpdate := s.repository.UpdateBalanceAccrued(ctx, order, &accrual)
				if errUpdate != nil {
					slog.Info("orders processing", "error", errUpdate.Error())
					done <- struct{}{}
					continue
				}

				done <- struct{}{}
			}
		}(jobs, done)
	}

	for _, orderNumber := range ordersToProcess {
		jobs <- orderNumber
	}
	close(jobs)

	for range len(ordersToProcess) {
		<-done
	}
}
