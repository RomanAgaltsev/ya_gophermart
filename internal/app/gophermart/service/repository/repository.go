package repository

import (
    "context"

    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/balance"
    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/order"
    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/user"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"

    "github.com/jackc/pgx/v5/pgxpool"
)

var _ user.Repository = (*Repo)(nil)
var _ order.Repository = (*Repo)(nil)
var _ balance.Repository = (*Repo)(nil)

func New(dbpool *pgxpool.Pool) *Repo {
    return &Repo{
        dbpool: dbpool,
    }
}

type Repo struct {
    dbpool *pgxpool.Pool
}

func (r *Repo) CreateUser(ctx context.Context, user *model.User) error {
    return nil
}

func (r *Repo) GetUser(ctx context.Context, login string) (*model.User, error) {
    return nil, nil
}

func (r *Repo) CreateOrder(ctx context.Context, order *model.Order) error {
    return nil
}

func (r *Repo) GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error) {
    return nil, nil
}
