package repository

import (
    "context"

    "github.com/RomanAgaltsev/ya_gophermart/internal/model"

    "github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
    Create(ctx context.Context, user model.User) error
}

func New(dbpool *pgxpool.Pool) *Repo {
    return &Repo{
        dbpool: dbpool,
    }
}

type Repo struct {
    dbpool *pgxpool.Pool
}

func (r *Repo) Create(ctx context.Context, user model.User) error {
    return nil
}
