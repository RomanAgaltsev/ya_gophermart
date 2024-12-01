package app

import (
	"context"
	"net/http"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/server"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/logger"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
)

type UserService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) error
}

type OrderService interface {
	Create(ctx context.Context, order *model.Order) error
	UserOrders(ctx context.Context, user *model.User) (model.Orders, error)
}

type BalanceService interface {
	UserBalance(ctx context.Context, user *model.User) (*model.Balance, error)
	BalanceWithdraw(ctx context.Context, user *model.User, order *model.Order, sum float64) error
	UserWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error)
}

// App struct of the application.
type App struct {
	cfg    *config.Config
	server *http.Server
}

// New creates new application.
func New() (*App, error) {
	app := &App{}

	// Configuration initialization
	err := app.initConfig()
	if err != nil {
		return nil, err
	}

	// Logger initialization
	err = app.initLogger()
	if err != nil {
		return nil, err
	}

	// HTTP server initialization
	err = app.initServer()
	if err != nil {
		return nil, err
	}

	return app, nil
}

// initConfig initializes application configuration.
func (a *App) initConfig() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	a.cfg = cfg

	return nil
}

// initLogger initializes logger.
func (a *App) initLogger() error {
	err := logger.Initialize()
	if err != nil {
		return err
	}

	return nil
}

// initServer initializes HTTP server.
func (a *App) initServer() error {
	srvr, err := server.New(a.cfg)
	if err != nil {
		return err
	}
	a.server = srvr

	return nil
}

// Run runs the application.
func (a *App) Run() error {
	return nil
}
