package repository

import (
    "context"

    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/balance"
    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/order"
    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/user"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"

    "github.com/jackc/pgx/v5/pgxpool"
)

var _ user.Repository = (*Repository)(nil)
var _ order.Repository = (*Repository)(nil)
var _ balance.Repository = (*Repository)(nil)

func New(dbpool *pgxpool.Pool) *Repository {
    return &Repository{
        dbpool: dbpool,
    }
}

type Repository struct {
    dbpool *pgxpool.Pool
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
    return nil
}

func (r *Repository) GetUser(ctx context.Context, login string) (*model.User, error) {
    return nil, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order *model.Order) error {
    return nil
}

func (r *Repository) GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error) {
    return nil, nil
}

func (r *Repository) GetBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
    return nil, nil
}

func (r *Repository) Withdraw(ctx context.Context, user *model.User, order *model.Order, sum float64) error {
    return nil
}

func (r *Repository) GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
    return nil, nil
}
