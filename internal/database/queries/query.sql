-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2) RETURNING id;

-- name: GetUser :one
SELECT id, login, password, created_at
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
SELECT id, login, number, status, accrual, uploaded_at
FROM orders
WHERE number = $1 LIMIT 1;

-- name: ListOrders :many
SELECT id, login, number, status, accrual, uploaded_at
FROM orders
WHERE login = $1
ORDER BY uploaded_at DESC;

-- name: ListOrdersToProcess :many
SELECT id, login, number, status, accrual, uploaded_at
FROM orders
WHERE status = 'NEW'
   OR status = 'PROCESSING';

-- name: CreateWithdraw :one
INSERT INTO withdrawals (login, order_number, sum)
VALUES ($1, $2, $3) RETURNING id;

-- name: ListWithdrawals :many
SELECT id, login, order_number, sum, processed_at
FROM withdrawals
WHERE login = $1
ORDER BY processed_at DESC;

-- name: CreateBalance :one
INSERT INTO balance (login)
VALUES ($1) RETURNING id;

-- name: GetBalance :one
SELECT id, login, accrued, withdrawn
FROM balance
WHERE login = $1 LIMIT 1;

-- name: UpdateBalanceAccrued :one
UPDATE balance
SET accrued = $2
WHERE login = $1 RETURNING accrued, withdrawn;

-- name: UpdateBalanceWithdrawn :one
UPDATE balance
SET withdrawn = withdrawn + $2
WHERE login = $1 RETURNING accrued, withdrawn;