-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2) RETURNING id;

-- name: GetUser :one
SELECT *
FROM users
WHERE login = $1 LIMIT 1;

-- name: CreateOrder :one
INSERT INTO orders (login, number)
VALUES ($1, $2) RETURNING id;

-- name: UpdateOrder :exec
UPDATE orders
SET status  = $2,
    accrual = $3
WHERE number = $1;

-- name: GetOrder :one
SELECT *
FROM orders
WHERE number = $1 LIMIT 1;

-- name: ListOrders :many
SELECT *
FROM orders
WHERE login = $1
ORDER BY uploaded_at DESC;

-- name: CreateWithdraw :one
INSERT INTO withdrawals (login, order_number, sum)
VALUES ($1, $2, $3) RETURNING id;

-- name: ListWithdrawals :many
SELECT *
FROM withdrawals
WHERE login = $1
ORDER BY processed_at DESC;

-- name: CreateBalance :one
INSERT INTO balance (login)
VALUES ($1) RETURNING id;

-- name: GetBalance :one
SELECT *
FROM balance
WHERE login = $1 LIMIT 1;

-- name: UpdateBalanceAccrued :exec
UPDATE balance
SET accrued = accrued + $2
WHERE login = $1 RETURNING accrued, withdrawn;

-- name: UpdateBalanceWithdrawn :exec
UPDATE balance
SET withdrawn = withdrawn + $2
WHERE login = $1 RETURNING accrued, withdrawn;