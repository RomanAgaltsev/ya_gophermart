-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    login      VARCHAR(20) UNIQUE NOT NULL,
    password   VARCHAR(20)        NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT NOW()
);

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY,
    number      VARCHAR(100) UNIQUE NOT NULL,
    status      order_status        NOT NULL DEFAULT 'NEW',
    accrual     NUMERIC(15, 3),
    uploaded_at TIMESTAMP           NOT NULL DEFAULT NOW()
);

CREATE TABLE withdrawals
(
    id           SERIAL PRIMARY KEY,
    order        VARCHAR(100)   NOT NULL,
    sum          NUMERIC(15, 3) NOT NULL,
    processed_at TIMESTAMP      NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd