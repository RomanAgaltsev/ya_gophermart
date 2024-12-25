package repository

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/RomanAgaltsev/ya_gophermart/internal/database/queries"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"

    "github.com/cenkalti/backoff/v4"
    "github.com/jackc/pgerrcode"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

var (
    ErrConflict        = fmt.Errorf("data conflict")
    ErrNegativeBalance = fmt.Errorf("negative balance")
)

// conflictOrder contains confict order and an error.
type conflictOrder struct {
    order *model.Order
    err   error
}

// PgxPool needs to mock pgxpool in tests.
type PgxPool interface {
    Close()
    Acquire(ctx context.Context) (c *pgxpool.Conn, err error)
    AcquireFunc(ctx context.Context, f func(*pgxpool.Conn) error) error
    AcquireAllIdle(ctx context.Context) []*pgxpool.Conn
    Reset()
    Config() *pgxpool.Config
    Stat() *pgxpool.Stat
    Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
    SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
    Begin(ctx context.Context) (pgx.Tx, error)
    BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
    CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
    Ping(ctx context.Context) error
}

// New creates new repository.
func New(dbpool PgxPool) (*Repository, error) {
    // Return Repository struct with new queries
    return &Repository{
        db: dbpool,
        q:  queries.New(dbpool),
    }, nil
}

// Repository is the repository structure.
type Repository struct {
    db PgxPool
    q  *queries.Queries
}

// CreateUser creates new user in the repository.
func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
    // PG error to catch the conflict
    var pgErr *pgconn.PgError

    // Create a function to wrap user creation with exponential backoff
    f := func() (error, error) {
        // Create user
        _, err := r.q.CreateUser(ctx, queries.CreateUserParams{
            Login:    user.Login,
            Password: user.Password,
        })

        // Check if there is a conflict
        if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
            return ErrConflict, nil
        }

        // Check if something has gone wrong
        if err != nil {
            return nil, err
        }

        return nil, nil
    }

    // Call the wrapping function
    errConf, err := backoff.RetryWithData(f, backoff.NewExponentialBackOff())
    if err != nil {
        return err
    }

    // There is a conflict
    if errConf != nil {
        return errConf
    }

    return nil
}

// GetUser returns a user from repository.
func (r *Repository) GetUser(ctx context.Context, login string) (*model.User, error) {
    // Get user from DB
    usr, err := backoff.RetryWithData(func() (queries.User, error) {
        return r.q.GetUser(ctx, login)
    }, backoff.NewExponentialBackOff())

    // Check if something has gone wrong
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return nil, err
    }

    // Check if there is nothing to return
    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }

    // Return user
    return &model.User{
        Login:    usr.Login,
        Password: usr.Password,
    }, nil
}

// CreateOrder creates new order in the repository.
func (r *Repository) CreateOrder(ctx context.Context, order *model.Order) (*model.Order, error) {
    // PG error to catch the conflict
    var pgErr *pgconn.PgError

    // Wrap order creation in a function
    f := func() (conflictOrder, error) {
        var co conflictOrder
        // Try to create an order
        _, errStore := r.q.CreateOrder(ctx, queries.CreateOrderParams{
            Login:  order.Login,
            Number: order.Number,
        })

        // Check if there is a conflict
        if errors.As(errStore, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
            orderByNumber, errGet := backoff.RetryWithData(func() (queries.Order, error) {
                // Return existing order
                return r.q.GetOrder(ctx, order.Number)
            }, backoff.NewExponentialBackOff())

            // Something has gone wrong
            if errGet != nil {
                return co, errGet
            }

            // Return conflict structure with existing order and an error
            return conflictOrder{
                order: &model.Order{
                    Login:      orderByNumber.Login,
                    Number:     orderByNumber.Number,
                    Status:     orderByNumber.Status,
                    Accrual:    orderByNumber.Accrual,
                    UploadedAt: orderByNumber.UploadedAt,
                },
                err: ErrConflict,
            }, nil
        }

        return co, errStore
    }

    // Call the wrapping function
    confOrder, err := backoff.RetryWithData(f, backoff.NewExponentialBackOff())
    if err != nil {
        return nil, err
    }

    // There is a conflict
    if errors.Is(confOrder.err, ErrConflict) {
        return &model.Order{
            Login:      confOrder.order.Login,
            Number:     confOrder.order.Number,
            Status:     confOrder.order.Status,
            Accrual:    confOrder.order.Accrual,
            UploadedAt: confOrder.order.UploadedAt,
        }, confOrder.err
    }

    return nil, nil
}

// GetListOfOrders returns a list of user orders.
func (r *Repository) GetListOfOrders(ctx context.Context, user *model.User) (model.Orders, error) {
    // Get orders from DB
    ordersQuery, err := backoff.RetryWithData(func() ([]queries.Order, error) {
        return r.q.ListOrders(ctx, user.Login)
    }, backoff.NewExponentialBackOff())
    if err != nil {
        return nil, err
    }

    // Fill the slice of orders to return
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

// CreateBalance creates user balance.
func (r *Repository) CreateBalance(ctx context.Context, user *model.User) error {
    // Create new balance in DB
    _, err := backoff.RetryWithData(func() (int32, error) {
        return r.q.CreateBalance(ctx, user.Login)
    }, backoff.NewExponentialBackOff())

    if err != nil {
        return err
    }

    return nil
}

// GetBalance returns user balance.
func (r *Repository) GetBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
    // Get user balance from DB.
    balanceQuery, err := backoff.RetryWithData(func() (queries.Balance, error) {
        return r.q.GetBalance(ctx, user.Login)
    }, backoff.NewExponentialBackOff())

    // Something has gone wrong
    if err != nil {
        return nil, err
    }

    // Return the balance
    return &model.Balance{
        Current:   balanceQuery.Accrued - balanceQuery.Withdrawn,
        Withdrawn: balanceQuery.Withdrawn,
    }, nil
}

// WithdrawFromBalance - withdraw the given sum from the user balance if its enough to withdraw.
func (r *Repository) WithdrawFromBalance(ctx context.Context, user *model.User, orderNumber string, sum float64) error {
    // Begin transaction
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    // Defer transaction rollback
    defer func() { _ = tx.Rollback(ctx) }()

    // Create query with transaction
    qtx := r.q.WithTx(tx)

    // Withdraw from balance
    withdrawnRow, err := backoff.RetryWithData(func() (queries.UpdateBalanceWithdrawnRow, error) {
        return qtx.UpdateBalanceWithdrawn(ctx, queries.UpdateBalanceWithdrawnParams{
            Login:     user.Login,
            Withdrawn: sum,
        })
    }, backoff.NewExponentialBackOff())
    if err != nil {
        // TODO
        _ = tx.Rollback(ctx)
        return err
    }

    // If the balance has become negative after withdrawal,
    // rollback the transaction and return the negative balance error
    if withdrawnRow.Accrued-withdrawnRow.Withdrawn < 0 {
        _ = tx.Rollback(ctx)
        return ErrNegativeBalance
    }

    // Balance enough to withdraw - create new withdrawal in DB
    _, err = backoff.RetryWithData(func() (int32, error) {
        return qtx.CreateWithdraw(ctx, queries.CreateWithdrawParams{
            Login:       user.Login,
            OrderNumber: orderNumber,
            Sum:         sum,
        })
    }, backoff.NewExponentialBackOff())
    if err != nil {
        // TODO
        _ = tx.Rollback(ctx)
        return err
    }

    return tx.Commit(ctx)
}

// GetListOfWithdrawals returns a list of withdrawals from the user balance.
func (r *Repository) GetListOfWithdrawals(ctx context.Context, user *model.User) (model.Withdrawals, error) {
    // Get withdrawals from DB
    withdrawalsQuery, err := backoff.RetryWithData(func() ([]queries.Withdrawal, error) {
        return r.q.ListWithdrawals(ctx, user.Login)
    }, backoff.NewExponentialBackOff())
    if err != nil {
        return nil, err
    }

    // Fill the slice of withdrawals to return
    withdrawals := make([]*model.Withdrawal, 0, len(withdrawalsQuery))
    for _, withdrawal := range withdrawalsQuery {
        withdrawals = append(withdrawals, &model.Withdrawal{
            Login:       withdrawal.Login,
            OrderNumber: withdrawal.OrderNumber,
            Sum:         withdrawal.Sum,
            ProcessedAt: withdrawal.ProcessedAt,
        })
    }

    return withdrawals, nil
}

// GetListOfOrdersToProcess returns the orders that need to be processed - NEW and PROCESSING statuses.
func (r *Repository) GetListOfOrdersToProcess(ctx context.Context) (model.Orders, error) {
    // Get orders from DB
    ordersQuery, err := backoff.RetryWithData(func() ([]queries.Order, error) {
        return r.q.ListOrdersToProcess(ctx)
    }, backoff.NewExponentialBackOff())
    if err != nil {
        return nil, err
    }

    // Fill the slice of orders to return
    ordersToProcess := make([]*model.Order, 0, len(ordersQuery))
    for _, order := range ordersQuery {
        ordersToProcess = append(ordersToProcess, &model.Order{
            Login:      order.Login,
            Number:     order.Number,
            Status:     order.Status,
            Accrual:    order.Accrual,
            UploadedAt: order.UploadedAt,
        })
    }

    return ordersToProcess, nil
}

// UpdateBalanceAccrued encreases user balance.
func (r *Repository) UpdateBalanceAccrued(ctx context.Context, order *model.Order, accrual *model.OrderAccrual) error {
    // Begin transaction
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    // Defer transaction rollback
    defer func() { _ = tx.Rollback(ctx) }()

    // Get query with transaction
    qtx := r.q.WithTx(tx)

    // Add the sum to the user balance in DB
    _, err = backoff.RetryWithData(func() (queries.UpdateBalanceAccruedRow, error) {
        return qtx.UpdateBalanceAccrued(ctx, queries.UpdateBalanceAccruedParams{
            Login:   order.Login,
            Accrued: accrual.Accrual,
        })
    }, backoff.NewExponentialBackOff())
    if err != nil {
        // TODO
        _ = tx.Rollback(ctx)
        return err
    }

    // Update order in DB 
    err = backoff.Retry(func() error {
        return qtx.UpdateOrder(ctx, queries.UpdateOrderParams{
            Number:  order.Number,
            Status:  accrual.Status,
            Accrual: accrual.Accrual,
        })
    }, backoff.NewExponentialBackOff())
    if err != nil {
        // TODO
        _ = tx.Rollback(ctx)
        return err
    }

    return tx.Commit(ctx)
}
