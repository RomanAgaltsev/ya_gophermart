package handler

import (
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"net/http"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

type Handler struct {
	cfg *config.Config
}

func NewHandlers(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) UserRegistrion(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {

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
