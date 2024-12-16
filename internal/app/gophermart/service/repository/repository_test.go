package repository_test

import (
	"context"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/mocks/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

const (
	createBalance = `-- name: CreateBalance :one
INSERT INTO balance (login)
VALUES ($1) RETURNING id
`
	createOrder = `-- name: CreateOrder :one
INSERT INTO orders (login, number)
VALUES ($1, $2) RETURNING id
`
	createUser = `-- name: CreateUser :one
INSERT INTO users (login, password)
VALUES ($1, $2) RETURNING id
`
	createWithdraw = `-- name: CreateWithdraw :one
INSERT INTO withdrawals (login, order_number, sum)
VALUES ($1, $2, $3) RETURNING id
`
	getBalance = `-- name: GetBalance :one
SELECT id, login, accrued, withdrawn
FROM balance
WHERE login = $1 LIMIT 1
`
	getOrder = `-- name: GetOrder :one
SELECT id, login, number, status, accrual, uploaded_at
FROM orders
WHERE number = $1 LIMIT 1
`
	getUser = `-- name: GetUser :one
SELECT id, login, password, created_at
FROM users
WHERE login = $1 LIMIT 1
`
	listOrders = `-- name: ListOrders :many
SELECT id, login, number, status, accrual, uploaded_at
FROM orders
WHERE login = $1
ORDER BY uploaded_at DESC
`
	listOrdersToProcess = `-- name: ListOrdersToProcess :many
SELECT id, login, number, status, accrual, uploaded_at
FROM orders
WHERE status = 'NEW'
   OR status = 'PROCESSING'
`
	listWithdrawals = `-- name: ListWithdrawals :many
SELECT id, login, order_number, sum, processed_at
FROM withdrawals
WHERE login = $1
ORDER BY processed_at DESC
`
	updateBalanceAccrued = `-- name: UpdateBalanceAccrued :one
UPDATE balance
SET accrued = $2
WHERE login = $1 RETURNING accrued, withdrawn
`
	updateBalanceWithdrawn = `-- name: UpdateBalanceWithdrawn :one
UPDATE balance
SET withdrawn = withdrawn + $2
WHERE login = $1 RETURNING accrued, withdrawn
`
	updateOrder = `-- name: UpdateOrder :exec
UPDATE orders
SET status  = $2,
    accrual = $3
WHERE number = $1
`
)

var _ = Describe("Repository", func() {
	var (
		err error

		ctx context.Context

		ctrl     *gomock.Controller
		mockPool *pgxpool.MockPgxPool
		repo     *repository.Repository

		user *model.User
	)

	BeforeEach(func() {
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		Expect(ctrl).ShouldNot(BeNil())

		mockPool = pgxpool.NewMockPgxPool(ctrl)

		repo, err = repository.New(mockPool)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Executing CreateUser method", func() {
		When("user doesn't exist", func() {
			BeforeEach(func() {
				login := "user"
				password := "password"
				//insertedID := 1

				user = &model.User{
					Login:    login,
					Password: password,
				}
				//pgxRow := NewRow(columns, 1, 2.3)
				mockPool.EXPECT().QueryRow(ctx, createUser, login, password).Return(nil)
			})

			It("returns nil error", func() {
				err = repo.CreateUser(ctx, user)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("user already exist", func() {
			BeforeEach(func() {
				login := "user"
				password := "password"
				//insertedID := 0

				user = &model.User{
					Login:    login,
					Password: password,
				}

				mockPool.EXPECT().QueryRow(ctx, createUser, login, password).Return(nil)
			})

			It("returns data conflict error", func() {
				err = repo.CreateUser(ctx, user)
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(pgconn.PgError{}))
			})
		})
	})
})
