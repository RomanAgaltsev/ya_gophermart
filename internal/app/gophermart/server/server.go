package server

import (
	"fmt"
	"net/http"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/api"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var ErrRunAddressIsEmpty = fmt.Errorf("configuration: HTTP server run address is empty")

// New creates new http server with middleware and routes
func New(cfg *config.Config) (*http.Server, error) {
	if cfg.RunAddress == "" {
		return nil, ErrRunAddressIsEmpty
	}

	// Create handler
	handle := api.New(cfg)

	// Create router
	router := chi.NewRouter()

	// Enable common middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5, api.ContentTypeJSON, api.ContentTypeText))

	// Set routes
	// -- public routes
	router.Group(func(r chi.Router) {
		r.Post("/api/user/register", handle.UserRegistrion)
		r.Post("/api/user/login", handle.UserLogin)
	})
	// -- protected routes
	router.Group(func(r chi.Router) {
		tokenAuth := jwtauth.New("HS256", []byte(cfg.SecretKey), nil)
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/api/user/orders", handle.OrderNumberUpload)
		r.Get("/api/user/orders", handle.OrderListRequest)
		r.Get("/api/user/balance", handle.UserBalanceRequest)
		r.Post("/api/user/balance/withdraw", handle.WithdrawalRequest)
		r.Get("/api/user/withdrawals", handle.WithdrawalsInformationRequest)
	})

	return &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}, nil
}
