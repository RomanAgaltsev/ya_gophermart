package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/RomanAgaltsev/ya_gophermart/internal/config"
    "github.com/RomanAgaltsev/ya_gophermart/internal/database"
    "github.com/RomanAgaltsev/ya_gophermart/internal/database/queries"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"

    "github.com/cenkalti/backoff/v4"
    "github.com/jackc/pgerrcode"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

var (
    ErrConflict = fmt.Errorf("data conflict")
)

type conflictOrder struct {
    order *model.Order
    err   error
}

func New(cfg *config.Config) (*Repository, error) {
    // Create context
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Create connection pool
    dbpool, err := database.NewConnectionPool(ctx, cfg.DatabaseURI)
    if err != nil {
        return nil, err
    }

    // Return Repository struct with new queries
    return &Repository{
        db: dbpool,
        q:  queries.New(dbpool),
    }, nil
}

type Repository struct {
    db *pgxpool.Pool
    q  *queries.Queries
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
    var pgErr *pgconn.PgError

    f := func() (error, error) {
        _, err := r.q.CreateUser(ctx, queries.CreateUserParams{
            Login:    user.Login,
            Password: user.Password,
        })

        if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
            return ErrConflict, nil
        }

        if err != nil {
            return nil, err
        }

        return nil, nil
    }

    errConf, err := backoff.RetryWithData(f, backoff.NewExponentialBackOff())
    if err != nil {
        return err
    }

    if errConf != nil {
        return errConf
    }

    return nil
}

func (r *Repository) GetUser(ctx context.Context, login string) (*model.User, error) {
    usr, err := backoff.RetryWithData(func() (queries.User, error) {
        return r.q.GetUser(ctx, login)
    }, backoff.NewExponentialBackOff())

    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return nil, err
    }

    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }

    return &model.User{
        Login:    usr.Login,
        Password: usr.Password,
    }, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
    var pgErr *pgconn.PgError

    f := func() (conflictOrder, error) {
        var co conflictOrder

        _, errStore := r.q.CreateOrder(ctx, queries.CreateOrderParams{
            Login:  order.Login,
            Number: order.Number,
        })

        if errors.As(errStore, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
            orderByNumber, errGet := backoff.RetryWithData(func() (queries.Order, error) {
                return r.q.GetOrder(ctx, order.Number)
            }, backoff.NewExponentialBackOff())

            if errGet != nil {
                return co, errGet
            }

            return conflictOrder{
                order: &model.Order{
                    Login:      orderByNumber.Login,
                    Number:     orderByNumber.Number,
                    Status:     orderByNumber.Status,
                    Accrual:    orderByNumber.Accrual,
                    UploadedAt: orderByNumber.UploadedAt,
                },
            }, nil
        }

        return co, errStore
    }

    confOrder, err := backoff.RetryWithData(f, backoff.NewExponentialBackOff())
    if err != nil {
        return nil, err
    }

    if errors.Is(confOrder.err, ErrConflict) {
        return &model.Order{
            Number:     confOrder.order.Number,
            Status:     confOrder.order.Status,
            Accrual:    confOrder.order.Accrual,
            UploadedAt: confOrder.order.UploadedAt,
        }, confOrder.err
    }

    return nil, nil
}

func (r *Repository) GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error) {
    ordersQuery, err := backoff.RetryWithData(func() ([]queries.Order, error) {
        return r.q.ListOrders(ctx, user.Login)
    }, backoff.NewExponentialBackOff())
    if err != nil {
        return nil, err
    }

    orders := make([]*model.Order, 0, len(ordersQuery))
    for _, order := range ordersQuery {
        orders = append(orders, &model.Order{
            Login:      order.Login,
            Number:     order.Number,
            Status:     order.Status,
            Accrual:    order.Accrual,
            UploadedAt: order.UploadedAt,
        })
    }

    return orders, nil
}

func (r *Repository) GetBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
    return nil, nil
}

func (r *Repository) Withdraw(ctx context.Context, user *model.User, orderNumber string, sum float64) error {
    return nil
}

func (r *Repository) GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
    withdrawalsQuery, err := backoff.RetryWithData(func() ([]queries.Withdrawal, error) {
        return r.q.
    }, backoff.NewExponentialBackOff())

    return nil, nil
}
