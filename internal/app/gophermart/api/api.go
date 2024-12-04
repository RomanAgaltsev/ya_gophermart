package api

import (
	"encoding/json"
	//"errors"
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

type Handler struct {
	cfg *config.Config

	userService    user.Service
	orderService   order.Service
	balanceService balance.Service
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) UserRegistrion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	var usr model.User
	if err := decoder.Decode(&usr); err != nil {
		slog.Info("decoding user", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusBadRequest)
		return
	}

	err := h.userService.Register(ctx, &usr)
	if err != nil {
		http.Error(w, "please look at logs", http.StatusInternalServerError)
	}
	//	if err != nil && !errors.Is(err, user.ErrLogin) {
	//		slog.Info("failed to short URL", "error", err.Error())
	//		http.Error(w, "please look at logs", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	if errors.Is(err, user.ErrConflict) {
	//		http.Error(w, "please look at logs", http.StatusConflict)
	//	}

	ja := auth.NewAuth(h.cfg.SecretKey)
	_, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
	if err != nil {
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	defer func() { _ = r.Body.Close() }()

	var usr model.User
	if err := decoder.Decode(&usr); err != nil {
		slog.Info("decoding user", "error", err.Error())
		http.Error(w, "please look at logs", http.StatusBadRequest)
		return
	}

	err := h.userService.Login(ctx, &usr)
	if err != nil {
		http.Error(w, "please look at logs", http.StatusInternalServerError)
	}

	ja := auth.NewAuth(h.cfg.SecretKey)
	_, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
	if err != nil {
		http.Error(w, "please look at logs", http.StatusInternalServerError)
		return
	}
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
