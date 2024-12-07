-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2) RETURNING id;

-- name: GetUser :one
SELECT *
FROM users
WHERE login = $1 LIMIT 1;

-- name: GetOrder :one
SELECT *
FROM orders
WHERE number = $1 LIMIT 1;

-- name: CreateOrder :one
INSERT INTO orders (login, number)
VALUES ($1, $2) RETURNING id;

-- name: ListOrders :many
SELECT *
FROM orders
WHERE login = $1
ORDER BY uploaded_at DESC;

-- name: CreateWithdraw :one
INSERT INTO withdrawals (order_number, sum)
VALUES ($1, $2) RETURNING id;

-- name: ListWithdrawals :many
SELECT *
FROM withdrawals
WHERE login = $1 LIMIT 1;
