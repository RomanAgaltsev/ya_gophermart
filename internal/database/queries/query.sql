-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2) RETURNING id;

-- name: GetUser :one
SELECT *
FROM users
WHERE login = $1 LIMIT 1;

-- name: ListOrders :many
SELECT *
FROM orders
WHERE user = $1
ORDER BY uploaded_at DESC;