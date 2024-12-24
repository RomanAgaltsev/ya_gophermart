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

// Service is the balance service interface.
type Service interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, user *model.User) (*model.Balance, error)
	Withdraw(ctx context.Context, user *model.User, orderNumber string, sum float64) error
	Withdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
}

// Repository is the balance service repository interface.
type Repository interface {
	CreateBalance(ctx context.Context, user *model.User) error
	GetBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	WithdrawFromBalance(ctx context.Context, user *model.User, orderNumber string, sum float64) error
	GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
	GetListOfOrdersToProcess(ctx context.Context) (model.Orders, error)
	UpdateBalanceAccrued(ctx context.Context, order *model.Order, accrual *model.OrderAccrual) error
}

// NewService creates new balance service.
func NewService(ctx context.Context, repository Repository, cfg *config.Config, runProcessing bool) (Service, error) {
	balanceService := &service{
		repository: repository,
		cfg:        cfg,
	}
	// Run orders processing goroutine only if needed
	if runProcessing {
		go balanceService.ordersProcessing(ctx)
	}

	return balanceService, nil
}

// service is the balance service structure.
type service struct {
	repository Repository
	cfg        *config.Config
}

// Create creates new user balance.
func (s *service) Create(ctx context.Context, user *model.User) error {
	return s.repository.CreateBalance(ctx, user)
}

// Get returns user balance.
func (s *service) Get(ctx context.Context, user *model.User) (*model.Balance, error) {
	return s.repository.GetBalance(ctx, user)
}

// Withdraw creates a withdrawal from user balance.
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

// Withdrawals returns a list of withdrawals from user balance.
func (s *service) Withdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
	return s.repository.GetListOfWithdrawals(ctx, user)
}

// ordersProcessing runs orders processing every 10 seconds.
func (s *service) ordersProcessing(ctx context.Context) {
	const ordersProcessingInterval = 10

	slog.Info("starting order processing")

	ticker := time.NewTicker(ordersProcessingInterval * time.Second)

	for {
		select {
		case <-ticker.C:
			slog.Info("order processing execution")
			s.processOrders()
		case <-ctx.Done():
			slog.Info("order processing stopped")
			return
		default:
			continue
		}
	}
}

// processOrders processes unprocessed orders.
func (s *service) processOrders() {
	// Fix the number of workers
	const workersNumber = 3

	// Create context
	ctx := context.Background()

	// Get orders to process - NEW and PROCESSING statuses
	ordersToProcess, err := s.repository.GetListOfOrdersToProcess(ctx)
	if err != nil {
		slog.Info("orders processing", "error", err.Error())
		return
	}

	// Create a channel for processing jobs
	jobs := make(chan *model.Order, len(ordersToProcess))
	// Create a channel for jobs completion awating
	done := make(chan struct{}, len(ordersToProcess))

	// Run workers
	for range workersNumber {
		// Every worker lives in its own goroutine
		go func(jobs chan *model.Order, done chan struct{}) {
			// Get orders from job channel
			for order := range jobs {
				// Get data from accrual system
				accrual, err := orderAccrual(s.cfg.AccrualSystemAddress, order.Number)
				if err != nil {
					slog.Info("orders processing", "error", err.Error())
					done <- struct{}{}
					continue
				}

				// If order status has not been changed, do nothing
				if order.Status == accrual.Status {
					done <- struct{}{}
					continue
				}

				// If order status has been changed, update balance
				errUpdate := s.repository.UpdateBalanceAccrued(ctx, order, accrual)
				if errUpdate != nil {
					slog.Info("orders processing", "error", errUpdate.Error())
					done <- struct{}{}
					continue
				}

				// Done with the order
				done <- struct{}{}
			}
		}(jobs, done)
	}

	// Fill the jobs channel with orders numbers
	for _, orderNumber := range ordersToProcess {
		jobs <- orderNumber
	}
	// Close jobs channel after filling
	close(jobs)

	// Waiting for all orders to be processed
	for range len(ordersToProcess) {
		<-done
	}
}

// orderAccrual fetches order data from the external accrual system.
func orderAccrual(accrualSystemAddress string, orderNumber string) (*model.OrderAccrual, error) {
	// Create HTTP client
	client := http.Client{}

	// Send request to the accrual system with exponential backoff
	resp, err := backoff.RetryWithData(func() (*http.Response, error) {
		url := fmt.Sprintf("%s/api/orders/%s", accrualSystemAddress, orderNumber)
		slog.Info("accrual system request", "address", url, "order", orderNumber)
		return client.Get(url)
	}, backoff.NewExponentialBackOff())
	defer func() { _ = resp.Body.Close() }()

	// Something has gone wrong
	if err != nil {
		return nil, err
	}

	// Decode order accrual data from JSON to accrual structure
	var accrual model.OrderAccrual
	err = render.DecodeJSON(resp.Body, &accrual)
	if err != nil {
		return nil, err
	}

	return &accrual, nil
}
