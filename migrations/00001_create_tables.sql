-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    login      VARCHAR(20) UNIQUE NOT NULL,
    password   VARCHAR(60)        NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT NOW()
);

CREATE TYPE order_status_type AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders
(
    id          SERIAL PRIMARY KEY,
    login       VARCHAR(20)         NOT NULL,
    number      VARCHAR(100) UNIQUE NOT NULL,
    status      order_status_type   NOT NULL DEFAULT 'NEW',
    accrual     DOUBLE PRECISION    NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMP           NOT NULL DEFAULT NOW()
);

CREATE TABLE withdrawals
(
    id           SERIAL PRIMARY KEY,
    login        VARCHAR(20)      NOT NULL,
    order_number VARCHAR(100)     NOT NULL,
    sum          DOUBLE PRECISION NOT NULL,
    processed_at TIMESTAMP        NOT NULL DEFAULT NOW()
);

CREATE TABLE balance
(
    id        SERIAL PRIMARY KEY,
    login     VARCHAR(20)      NOT NULL,
    accrued   DOUBLE PRECISION NOT NULL DEFAULT 0,
    withdrawn DOUBLE PRECISION NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE balance;
DROP TABLE withdrawals;
DROP TABLE orders;
DROP TYPE order_status;
DROP TABLE users;
-- +goose StatementEnd