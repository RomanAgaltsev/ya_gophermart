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
    contentTypeJSON = "application/json"

    argError = "error"

    msgNewJWTToken       = "new JWT token"
    msgUserRegistration  = "user registration"
    msgUserLogin         = "user login"
    msgOrderNumberUpload = "order number upload"
    msgOrderList         = "get orders list"
    msgNewUserBalance    = "new user balance"
    msgUserBalance       = "user balance request"
    msgWithdraw          = "withdraw request"
    msgUserWithdrawals   = "user withdrawals request"
)

// Handler handles all HTTP requests.
type Handler struct {
    cfg *config.Config

    userService    user.Service
    orderService   order.Service
    balanceService balance.Service
}

// NewHandler is a Handler constructor.
func NewHandler(cfg *config.Config, userService user.Service, orderService order.Service, balanceService balance.Service) *Handler {
    return &Handler{
        cfg:            cfg,
        userService:    userService,
        orderService:   orderService,
        balanceService: balanceService,
    }
}

// UserRegistrion handles user registration request.
func (h *Handler) UserRegistrion(w http.ResponseWriter, r *http.Request) {
    var usr model.User
    if err := render.Bind(r, &usr); err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }

    // Get context from request
    ctx := r.Context()

    // Register user
    err := h.userService.Register(ctx, &usr)
    if err != nil && !errors.Is(err, user.ErrLoginIsAlreadyTaken) {
        // There is an error, but not a conflict
        slog.Info(msgUserRegistration, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if errors.Is(err, user.ErrLoginIsAlreadyTaken) {
        // There is a conflict
        slog.Info(msgUserRegistration, argError, err.Error())
        _ = render.Render(w, r, ErrLoginIsAlreadyTaken)
        return
    }

    // Create a balance for the user
    err = h.balanceService.Create(ctx, &usr)
    if err != nil {
        slog.Info(msgNewUserBalance, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    // Generate JWT token
    ja := auth.NewAuth(h.cfg.SecretKey)
    _, tokenString, err := auth.NewJWTToken(ja, usr.Login)
    if err != nil {
        // Something has gone wrong
        slog.Info(msgNewJWTToken, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    // Set a cookie with generated JWT token
    http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))

    w.WriteHeader(http.StatusOK)
}

// UserLogin handles user login request.
func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
    var usr model.User
    if err := render.Bind(r, &usr); err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }

    // Get context from request
    ctx := r.Context()

    // Login user
    err := h.userService.Login(ctx, &usr)
    if err != nil && !errors.Is(err, user.ErrWrongLoginPassword) {
        // There is an error, but not with login/password pair
        slog.Info(msgUserLogin, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if errors.Is(err, user.ErrWrongLoginPassword) {
        // There is a problem with login/password
        slog.Info(msgUserLogin, argError, err.Error())
        _ = render.Render(w, r, ErrWrongLoginPassword)
        return
    }

    // Generate JWT token
    ja := auth.NewAuth(h.cfg.SecretKey)
    _, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
    if err != nil {
        // Something has gone wrong
        slog.Info(msgNewJWTToken, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    // Set a cookie with generated JWT token
    http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) OrderNumberUpload(w http.ResponseWriter, r *http.Request) {
    rBody, _ := io.ReadAll(r.Body)
    defer func() { _ = r.Body.Close() }()

    orderNumber := string(rBody)

    if orderNumber == "" {
        _ = render.Render(w, r, ErrBadRequest)
        return
    }

    if !orderpkg.IsNumberValid(orderNumber) {
        _ = render.Render(w, r, ErrInvalidOrderNumber)
        return
    }

    // Get context from request
    ctx := r.Context()

    usr, err := auth.UserFromRequest(r, h.cfg.SecretKey)
    if err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }

    ordr := model.Order{
        Login:  usr.Login,
        Number: orderNumber,
    }

    err = h.orderService.Create(ctx, &ordr)
    if err != nil && !errors.Is(err, order.ErrOrderUploadedByThisLogin) && !errors.Is(err, order.ErrOrderUploadedByAnotherLogin) {
        // There is an error, but not a conflict
        slog.Info(msgOrderNumberUpload, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if errors.Is(err, order.ErrOrderUploadedByThisLogin) {
        // There is a conflict
        slog.Info(msgOrderNumberUpload, argError, err.Error())
        _ = render.Render(w, r, ErrOrderUploadedByThisLogin)
        return
    }

    if errors.Is(err, order.ErrOrderUploadedByAnotherLogin) {
        // There is a conflict
        slog.Info(msgOrderNumberUpload, argError, err.Error())
        _ = render.Render(w, r, ErrOrderUploadedByAnotherLogin)
        return
    }

    w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) OrderListRequest(w http.ResponseWriter, r *http.Request) {
    // Get context from request
    ctx := r.Context()

    usr, err := auth.UserFromRequest(r, h.cfg.SecretKey)
    if err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }

    orders, err := h.orderService.UserOrders(ctx, usr)
    if err != nil {
        slog.Info(msgOrderList, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if len(orders) == 0 {
        _ = render.Render(w, r, ErrNoOrders)
        return
    }

    w.Header().Set("Content-type", contentTypeJSON)
    w.WriteHeader(http.StatusOK)

    if err := render.Render(w, r, orders); err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }
}

func (h *Handler) UserBalanceRequest(w http.ResponseWriter, r *http.Request) {
    // Get context from request
    ctx := r.Context()

    usr, err := auth.UserFromRequest(r, h.cfg.SecretKey)
    if err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }

    userBalance, err := h.balanceService.Get(ctx, usr)
    if err != nil {
        slog.Info(msgUserBalance, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    w.Header().Set("Content-type", contentTypeJSON)
    w.WriteHeader(http.StatusOK)

    if err := render.Render(w, r, userBalance); err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }
}

func (h *Handler) WithdrawRequest(w http.ResponseWriter, r *http.Request) {
    var withdrawal model.Withdrawal

    if err := render.Bind(r, &withdrawal); err != nil {
        _ = render.Render(w, r, ErrBadRequest)
        return
    }

    if !orderpkg.IsNumberValid(withdrawal.OrderNumber) {
        _ = render.Render(w, r, ErrInvalidOrderNumber)
        return
    }

    // Get context from request
    ctx := r.Context()

    usr, err := auth.UserFromRequest(r, h.cfg.SecretKey)
    if err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }

    err = h.balanceService.Withdraw(ctx, usr, withdrawal.OrderNumber, withdrawal.Sum)
    if err != nil && !errors.Is(err, balance.ErrNotEnoughBalance) {
        // There is an error, but not with balance
        slog.Info(msgWithdraw, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if errors.Is(err, balance.ErrNotEnoughBalance) {
        // There is a problem with balance - not enough to withdraw the sum
        slog.Info(msgWithdraw, argError, err.Error())
        _ = render.Render(w, r, ErrNotEnoughBalance)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) WithdrawalsInformationRequest(w http.ResponseWriter, r *http.Request) {
    // Get context from request
    ctx := r.Context()

    usr, err := auth.UserFromRequest(r, h.cfg.SecretKey)
    if err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }

    withdrawals, err := h.balanceService.Withdrawals(ctx, usr)
    if err != nil {
        // There is an error, but not with withdrawals
        slog.Info(msgUserWithdrawals, argError, err.Error())
        _ = render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    if len(withdrawals) == 0 {
        _ = render.Render(w, r, ErrNoWithdrawals)
        return
    }

    w.Header().Set("Content-type", contentTypeJSON)
    w.WriteHeader(http.StatusOK)

    if err := render.Render(w, r, withdrawals); err != nil {
        _ = render.Render(w, r, ErrorRenderer(err))
        return
    }
}
