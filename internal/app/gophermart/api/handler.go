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

    var usr model.User
    if err := render.Bind(r, &usr); err != nil {
        render.Render(w, r, ErrBadRequest)
        return
    }

    // Register user
    err := h.userService.Register(ctx, &usr)
    if err != nil && !errors.Is(err, user.ErrLoginIsAlreadyTaken) {
        // There is an error, but not a conflict
        slog.Info("user registration", "error", err.Error())
        render.Render(w, r, ErrorRenderer(err))
        return
    }

    if errors.Is(err, user.ErrLoginIsAlreadyTaken) {
        // There is a conflict
        slog.Info("user registration", "error", err.Error())
        render.Render(w, r, ErrLoginIsAlreadyTaken)
    }

    // Generate JWT token
    ja := auth.NewAuth(h.cfg.SecretKey)
    _, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
    if err != nil {
        // Something has gone wrong
        slog.Info("new JWT token", "error", err.Error())
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
        slog.Info("user login", "error", err.Error())
        render.Render(w, r, ErrorRenderer(err))
        return
    }

    if errors.Is(err, user.ErrWrongLoginPassword) {
        // There is a problem with login/password
        slog.Info("user registration", "error", err.Error())
        render.Render(w, r, ErrWrongLoginPassword)
    }

    // Generate JWT token
    ja := auth.NewAuth(h.cfg.SecretKey)
    _, tokenString, _ := auth.NewJWTToken(ja, usr.Login)
    if err != nil {
        // Something has gone wrong
        slog.Info("new JWT token", "error", err.Error())
        render.Render(w, r, ServerErrorRenderer(err))
        return
    }

    // Set a cookie with generated JWT token
    http.SetCookie(w, auth.NewCookieWithDefaults(tokenString))

    w.WriteHeader(http.StatusOK)
}

func (h *Handler) OrderNumberUpload(w http.ResponseWriter, r *http.Request) {
    /*
       500 — внутренняя ошибка сервера.
    */

    // Get context from request
    ctx := r.Context()

    rBody, _ := io.ReadAll(r.Body)
    defer func() { _ = r.Body.Close() }()

    orderNumber := string(rBody)

    //400 — неверный формат запроса
    if orderNumber == "" {
        render.Render(w, r, ErrBadRequest)
        return
    }

    // 422 — неверный формат номера заказа
    if !orderpkg.IsNumberValid(orderNumber) {
        render.Render(w, r, ErrInvalidOrderNumber)
        return
    }

    //ordr := model.Order{}

    // 409 — номер заказа уже был загружен другим пользователем
    //    err := h.orderService.Create(ctx, &ordr)
    //    if err != nil && !errors.Is(err, user.ErrLoginIsAlreadyTaken) {
    //        // There is an error, but not a conflict
    //        slog.Info("user registration", "error", err.Error())
    //        render.Render(w, r, ErrorRenderer(err))
    //        return
    //    }

    //    if errors.Is(err, user.ErrLoginIsAlreadyTaken) {
    //        // There is a conflict
    //        slog.Info("user registration", "error", err.Error())
    //        render.Render(w, r, ErrLoginIsAlreadyTaken)
    //    }

    // 200 — номер заказа уже был загружен этим пользователем;
    w.WriteHeader(http.StatusOK)

    // 202 — новый номер заказа принят в обработку;
    w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) OrderListRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UserBalanceRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) WithdrawalRequest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) WithdrawalsInformationRequest(w http.ResponseWriter, r *http.Request) {

}
