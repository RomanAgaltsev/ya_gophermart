package app

import (
	"fmt"
	"net/http"

	handler "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/api/http"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var ErrRunAddressIsEmpty = fmt.Errorf("configuration: HTTP server run address is empty")

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
	// Check HTTP server run address
	if a.cfg.RunAddress == "" {
		return ErrRunAddressIsEmpty
	}

	// Create handler
	handle := handler.New(a.cfg)

	// Create router
	router := chi.NewRouter()

	// Enable common middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5, handler.ContentTypeJSON, handler.ContentTypeText))

	// Set routes
	// -- public routes
	router.Group(func(r chi.Router) {
		r.Post("/api/user/register", handle.UserRegistrion)
		r.Post("/api/user/login", handle.UserLogin)
	})
	// -- protected routes
	router.Group(func(r chi.Router) {
		tokenAuth := jwtauth.New("HS256", []byte(a.cfg.SecretKey), nil)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/api/user/orders", handle.OrderNumberUpload)
		r.Get("/api/user/orders", handle.OrderListRequest)
		r.Get("/api/user/balance", handle.UserBalanceRequest)
		r.Post("/api/user/balance/withdraw", handle.WithdrawalRequest)
		r.Get("/api/user/withdrawals", handle.WithdrawalsInformationRequest)
	})

	a.server = &http.Server{
		Addr:    a.cfg.RunAddress,
		Handler: router,
	}

	return nil
}

// Run runs the application.
func (a *App) Run() error {
	return nil
}
