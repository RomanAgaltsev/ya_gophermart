package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/balance"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/order"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/user"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
	"github.com/RomanAgaltsev/ya_gophermart/internal/pkg/auth"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

// Handler handles all HTTP requests.
type Handler struct {
	cfg *config.Config

	userService    user.Service
	orderService   order.Service
	balanceService balance.Service
}

// NewHandler is a Handler constructor.
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

// UserRegistrion handles user registration request.
func (h *Handler) UserRegistrion(w http.ResponseWriter, r *http.Request) {
	// Get context from request
	ctx := r.Context()

	// Create decoder
	decoder := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Decode user struct from request body
	var usr model.User
	if err := decoder.Decode(&usr); err != nil {
		// Something has gone wrong
		slog.Info("decoding user", "error", err.Error())
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	// Register user
	err := h.userService.Register(ctx, &usr)
	if err != nil && !errors.Is(err, user.ErrLoginIsAlreadyTaken) {
		// There is an error, but not a conflict
		slog.Info("user registration", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, user.ErrLoginIsAlreadyTaken) {
		// There is a conflict
		http.Error(w, user.ErrLoginIsAlreadyTaken.Error(), http.StatusConflict)
	}

	// Generate JWT token
	ja := auth.NewAuth(h.cfg.SecretKey)
	_, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
	if err != nil {
		// Something has gone wrong
		slog.Info("new JWT token", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Set a cookie with generated JWT token
	http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))

	w.WriteHeader(http.StatusOK)
}

// UserLogin handles user login request.
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	// Get context from request
	ctx := r.Context()

	// Create decoder
	decoder := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	// Decode user struct from request body
	var usr model.User
	if err := decoder.Decode(&usr); err != nil {
		// Something has gone wrong
		slog.Info("decoding user", "error", err.Error())
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	// Login user
	err := h.userService.Login(ctx, &usr)
	if err != nil && !errors.Is(err, user.ErrWrongLoginPassword) {
		// There is an error, but not with login/password pair
		slog.Info("user login", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	if errors.Is(err, user.ErrWrongLoginPassword) {
		// There is a problem with login/password
		http.Error(w, user.ErrWrongLoginPassword.Error(), http.StatusUnauthorized)
	}

	// Generate JWT token
	ja := auth.NewAuth(h.cfg.SecretKey)
	_, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
	if err != nil {
		// Something has gone wrong
		slog.Info("new JWT token", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}

	// Set a cookie with generated JWT token
	http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) OrderNumberUpload(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) OrderListRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UserBalanceRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) WithdrawalRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) WithdrawalsInformationRequest(w http.ResponseWriter, r *http.Request) {

}
