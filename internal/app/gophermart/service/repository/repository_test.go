package repository_test

import (
	"context"
	//"database/sql"
	"time"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/mocks/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

		user       *model.User
		expectUser *model.User
		login      string
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

	Context("Calling CreateUser method", func() {
		When("user doesn't exist", func() {
			BeforeEach(func() {
				login := "user"
				password := "password"
				var insertedID int32 = 1

				user = &model.User{
					Login:    login,
					Password: password,
				}
				pgxRow := pgxpool.NewRow(insertedID).WithError(nil)
				mockPool.EXPECT().QueryRow(ctx, createUser, login, password).Return(pgxRow)
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
				var insertedID int32 = 0

				user = &model.User{
					Login:    login,
					Password: password,
				}

				pgxRow := pgxpool.NewRow(insertedID).WithError(&pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation})
				mockPool.EXPECT().QueryRow(ctx, createUser, login, password).Return(pgxRow)
			})

			It("returns data conflict error", func() {
				err = repo.CreateUser(ctx, user)
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(repository.ErrConflict))
			})
		})
	})

	Context("Calling GetUser method", func() {
		// TODO
		XWhen("user doesn't exist", func() {
			BeforeEach(func() {
				var ID int32 = 0
				login = ""
				password := ""
				createdAt := time.Now()

				pgxRow := pgxpool.NewRow(ID, login, password, createdAt)
				//.WithError(sql.ErrNoRows)
				mockPool.EXPECT().QueryRow(ctx, getUser, login).Return(pgxRow)
			})

			It("returns nil user and nil error", func() {
				_, err := repo.GetUser(ctx, login)
				//Expect(err).Should(HaveOccurred())
				Expect(err).ShouldNot(HaveOccurred())
				//Expect(user).To(BeNil())
			})
		})

		When("user exists", func() {
			BeforeEach(func() {
				var ID int32 = 1
				login = "user"
				password := "password"
				createdAt := time.Now()

				expectUser = &model.User{
					Login:    login,
					Password: password,
				}

				pgxRow := pgxpool.NewRow(ID, login, password, createdAt)
				mockPool.EXPECT().QueryRow(ctx, getUser, login).Return(pgxRow)
			})

			It("returns a user and nil error", func() {
				user, err = repo.GetUser(ctx, login)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(*user).To(Equal(*expectUser))
			})
		})

		When("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil user and an error", func() {

			})
		})
	})

	Context("Calling CreateOrder method", func() {
		When("order doesn't exist", func() {
			BeforeEach(func() {

			})

			It("returns nil order and nil error", func() {

			})
		})

		When("order exists", func() {
			BeforeEach(func() {

			})

			It("returns the order and data conflict error", func() {

			})
		})

		When("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil order and an error", func() {

			})
		})
	})

	Context("Calling GetListOfOrders method", func() {
		When("orders exist", func() {
			BeforeEach(func() {

			})

			It("returns a non-empty list of orders and nil error", func() {

			})
		})

		When("orders don't exist", func() {
			BeforeEach(func() {

			})

			It("returns an empty list of orders and nil error", func() {

			})
		})

		When("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil list of orders and an error", func() {

			})
		})
	})

	Context("Calling CreateBalance method", func() {
		When("the balance is successfully created", func() {
			BeforeEach(func() {

			})

			It("returns nil error", func() {

			})
		})

		When("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns an error", func() {

			})
		})
	})

	Context("Calling GetBalance method", func() {
		When("there is no error", func() {
			BeforeEach(func() {

			})

			It("returns a balance and nil error", func() {

			})
		})

		When("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil balance and an error", func() {

			})
		})
	})

	Context("Calling WithdrawFromBalance method", func() {
		When("everything is right", func() {
			BeforeEach(func() {

			})

			It("returns nil error", func() {

			})
		})

		When("balance not enough to withdraw", func() {
			BeforeEach(func() {

			})

			It("returns negative balance error", func() {

			})
		})

		When("something has gone wrong with the queries", func() {
			BeforeEach(func() {

			})

			It("returns an error", func() {

			})
		})
	})

	Context("Calling GetListOfWithdrawals method", func() {
		When("withdrawals exist", func() {
			BeforeEach(func() {

			})

			It("returns a non-empty list of withdrawals and nil error", func() {

			})
		})

		When("withdrawals don't exist", func() {
			BeforeEach(func() {

			})

			It("returns an empty list of withdrawals and nil error", func() {

			})
		})

		When("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil list of withdrawals and an error", func() {

			})
		})
	})

	Context("Calling GetListOfOrdersToProcess method", func() {
		When("orders to process exist", func() {
			BeforeEach(func() {

			})

			It("returns a non-empty list of orders and nil error", func() {

			})
		})

		When("orders to process don't exist", func() {
			BeforeEach(func() {

			})

			It("returns an empty list of orders and nil error", func() {

			})
		})

		When("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil list of orders and an error", func() {

			})
		})
	})

	Context("Calling UpdateBalanceAccrued method", func() {
		When("everything is right", func() {
			BeforeEach(func() {

			})

			It("returns nil error", func() {

			})
		})

		When("something has gone wrong with the queries", func() {
			BeforeEach(func() {

			})

			It("returns an error", func() {

			})
		})
	})
})
