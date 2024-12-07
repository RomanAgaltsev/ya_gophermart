package api

import (
    "errors"
    "io"
    "log/slog"
    "net/http"

    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/balance"
    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/order"
    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/user"
    "github.com/RomanAgaltsev/ya_gophermart/internal/config"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"
    "github.com/RomanAgaltsev/ya_gophermart/internal/pkg/auth"
    orderpkg "github.com/RomanAgaltsev/ya_gophermart/internal/pkg/order"

    "github.com/go-chi/render"
)

const (
    argError = "error"

    msgNewJWTToken       = "new JWT token"
    msgUserRegistration  = "user registration"
    msgUserLogin         = "user login"
    msgOrderNumberUpload = "order number upload"
    msgOrderList         = "get orders list"
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

    var usr model.User
    if err := render.Bind(r, &usr); err != nil {
        render.Render(w, r, ErrBadRequest)
        return
    }

    // Register user
    err := h.userService.Register(ctx, &usr)
    if err != nil && !errors.Is(err, user.ErrLoginIsAlreadyTaken) {
        // There is an error, but not a conflict
        slog.Info(msgUserRegistration, argError, err.Error())
        render.Render(w, r, ErrorRenderer(err))
        return
    }

    if errors.Is(err, user.ErrLoginIsAlreadyTaken) {
        // There is a conflict
        slog.Info(msgUserRegistration, argError, err.Error())
        render.Render(w, r, ErrLoginIsAlreadyTaken)
    }

    // Generate JWT token
    ja := auth.NewAuth(h.cfg.SecretKey)
    _, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
    if err != nil {
        // Something has gone wrong
        slog.Info(msgNewJWTToken, argError, err.Error())
        render.Render(w, r, ServerErrorRenderer(err))
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

    var usr model.User
    if err := render.Bind(r, &usr); err != nil {
        render.Render(w, r, ErrBadRequest)
        return
    }

    // Login user
    err := h.userService.Login(ctx, &usr)
    if err != nil && !errors.Is(err, user.ErrWrongLoginPassword) {
        // There is an error, but not with login/password pair
        slog.Info(msgUserLogin, argError, err.Error())
        render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if errors.Is(err, user.ErrWrongLoginPassword) {
        // There is a problem with login/password
        slog.Info(msgUserLogin, argError, err.Error())
        render.Render(w, r, ErrWrongLoginPassword)
    }

    // Generate JWT token
    ja := auth.NewAuth(h.cfg.SecretKey)
    _, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
    if err != nil {
        // Something has gone wrong
        slog.Info(msgNewJWTToken, argError, err.Error())
        render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    // Set a cookie with generated JWT token
    http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) OrderNumberUpload(w http.ResponseWriter, r *http.Request) {
    // Get context from request
    ctx := r.Context()

    rBody, _ := io.ReadAll(r.Body)
    defer func() { _ = r.Body.Close() }()

    orderNumber := string(rBody)

    if orderNumber == "" {
        render.Render(w, r, ErrBadRequest)
        return
    }

    if !orderpkg.IsNumberValid(orderNumber) {
        render.Render(w, r, ErrInvalidOrderNumber)
        return
    }

    ordr := model.Order{
        Login:  "",
        Number: orderNumber,
    }

    err := h.orderService.Create(ctx, &ordr)
    if err != nil && !errors.Is(err, order.ErrOrderUploadedByThisLogin) && !errors.Is(err, order.ErrOrderUploadedByAnotherLogin) {
        // There is an error, but not a conflict
        slog.Info(msgOrderNumberUpload, argError, err.Error())
        render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if errors.Is(err, order.ErrOrderUploadedByThisLogin) {
        // There is a conflict
        slog.Info(msgOrderNumberUpload, argError, err.Error())
        render.Render(w, r, ErrOrderUploadedByThisLogin)
        return
    }

    if errors.Is(err, order.ErrOrderUploadedByAnotherLogin) {
        // There is a conflict
        slog.Info(msgOrderNumberUpload, argError, err.Error())
        render.Render(w, r, ErrOrderUploadedByAnotherLogin)

    }

    w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) OrderListRequest(w http.ResponseWriter, r *http.Request) {
    // Get context from request
    ctx := r.Context()

    usr := model.User{}

    orders, err := h.orderService.UserOrders(ctx, &usr)
    if err != nil {
        slog.Info(msgOrderList, argError, err.Error())
        render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if len(orders) == 0 {
        render.Render(w, r, ErrNoOrders)
        return
    }

    w.WriteHeader(http.StatusOK)

    if err := render.Render(w, r, orders); err != nil {
        render.Render(w, r, ErrorRenderer(err))
    }
}

func (h *Handler) UserBalanceRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) WithdrawalRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) WithdrawalsInformationRequest(w http.ResponseWriter, r *http.Request) {

}
