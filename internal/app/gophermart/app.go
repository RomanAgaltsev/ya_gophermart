package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/server"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/balance"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/order"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/user"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/logger"
)

// App struct of the application.
type App struct {
	cfg    *config.Config
	server *http.Server

	userService    user.Service
	orderService   order.Service
	balanceService balance.Service
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

	// Services initialization
	err = app.initServices()
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

// initServices initializes application services.
func (a *App) initServices() error {
	// Create repository
	repo, err := repository.New(a.cfg)
	if err != nil {
		return err
	}

	// Create user service
	userService, err := user.NewService(repo, a.cfg)
	if err != nil {
		return nil
	}
	a.userService = userService

	// Create order service
	orderService, err := order.NewService(repo, a.cfg)
	if err != nil {
		return nil
	}
	a.orderService = orderService

	// Create balance service
	balanceService, err := balance.NewService(repo, a.cfg, true)
	if err != nil {
		return nil
	}
	a.balanceService = balanceService

	return nil
}

// initServer initializes HTTP server.
func (a *App) initServer() error {
	srvr, err := server.New(a.cfg, a.userService, a.orderService, a.balanceService)
	if err != nil {
		return err
	}
	a.server = srvr

	return nil
}

// Run runs the application.
func (a *App) Run() error {
	return a.runApplication()
}

func (a *App) runApplication() error {
	// Create channels for graceful shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	// Interrupt signal
	signal.Notify(quit, os.Interrupt)

	// Graceful shutdown executes in a goroutine
	go func() {
		<-quit
		slog.Info("shutting down HTTP server")

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Shutdown HTTP server
		if err := a.server.Shutdown(ctx); err != nil {
			slog.Error("HTTP server shutdown error", slog.String("error", err.Error()),
			)
		}

		close(done)
	}()

	slog.Info("starting HTTP server", "addr", a.server.Addr)

	// Run HTTP server
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server error", slog.String("error", err.Error()))
		return err
	}

	<-done
	slog.Info("HTTP server stopped")
	return nil
}
