package repository_test

import (
	"context"
	"errors"
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
		err                 error
		errSomethingStrange error

		ctx context.Context

		ctrl     *gomock.Controller
		mockPool *pgxpool.MockPgxPool
		repo     *repository.Repository

		rowID int32

		// User
		userLogin     string
		userPassword  string
		userCreatedAt time.Time

		user         model.User
		userExpected model.User

		// Order
		orderNumber     string
		orderUploadedAt time.Time

		order         model.Order
		orderExpected model.Order
	)

	BeforeEach(func() {
		errSomethingStrange = errors.New("something strange")

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
				rowID = 1
				userLogin = "user"
				userPassword = "password"

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}
				pgxRow := pgxpool.NewRow(rowID).WithError(nil)
				mockPool.EXPECT().QueryRow(ctx, createUser, userLogin, userPassword).Return(pgxRow).Times(1)
			})

			It("returns nil error", func() {
				err = repo.CreateUser(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("user already exist", func() {
			BeforeEach(func() {
				rowID = 0
				userLogin = "user"
				userPassword = "password"

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				pgxRow := pgxpool.NewRow(rowID).WithError(&pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation})
				mockPool.EXPECT().QueryRow(ctx, createUser, userLogin, userPassword).Return(pgxRow).Times(1)
			})

			It("returns data conflict error", func() {
				err = repo.CreateUser(ctx, &user)
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(repository.ErrConflict))
			})
		})
	})

	Context("Calling GetUser method", func() {
		When("user doesn't exist", func() {
			BeforeEach(func() {
				rowID = 0
				userLogin = ""
				userPassword = ""
				userCreatedAt = time.Now()

				pgxRow := pgxpool.NewRow(rowID, userLogin, userPassword, userCreatedAt)
				mockPool.EXPECT().QueryRow(ctx, getUser, userLogin).Return(pgxRow).Times(1)
			})

			It("returns empty user and nil error", func() {
				result, err := repo.GetUser(ctx, userLogin)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Login).To(BeEmpty())
				Expect(result.Password).To(BeEmpty())
			})
		})

		When("user exists", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				userPassword = "password"
				userCreatedAt = time.Now()

				userExpected = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				pgxRow := pgxpool.NewRow(rowID, userLogin, userPassword, userCreatedAt)
				mockPool.EXPECT().QueryRow(ctx, getUser, userLogin).Return(pgxRow).Times(1)
			})

			It("returns a user and nil error", func() {
				result, err := repo.GetUser(ctx, userLogin)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(*result).To(Equal(userExpected))
			})
		})

		// Pending - because of exponential backoff it lasts about 10 seconds
		// Remove 'X' to run
		XWhen("something has gone wrong with the query", func() {
			BeforeEach(func() {
				rowID = 0
				userLogin = ""
				userPassword = ""
				userCreatedAt = time.Now()

				pgxRow := pgxpool.NewRow(rowID, userLogin, userPassword, userCreatedAt).WithError(errSomethingStrange)
				mockPool.EXPECT().QueryRow(ctx, getUser, userLogin).Return(pgxRow).AnyTimes()
			})

			It("returns nil user and an error", func() {
				result, err := repo.GetUser(ctx, userLogin)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
			})
		})
	})

	Context("Calling CreateOrder method", func() {
		When("order doesn't exist", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				orderNumber = "12345678903"

				order = model.Order{
					Login:  userLogin,
					Number: orderNumber,
				}

				pgxRow := pgxpool.NewRow(rowID)
				mockPool.EXPECT().QueryRow(ctx, createOrder, userLogin, orderNumber).Return(pgxRow).Times(1)
			})

			It("returns nil order and nil error", func() {
				result, err := repo.CreateOrder(ctx, &order)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).To(BeNil())
			})
		})

		When("order exists", func() {
			BeforeEach(func() {
				rowID = 0
				userLogin = "user"
				orderNumber = "12345678903"
				orderUploadedAt = time.Now()

				order = model.Order{
					Login:  userLogin,
					Number: orderNumber,
				}
				orderExpected = model.Order{
					Login:      "another user",
					Number:     orderNumber,
					Status:     "NEW",
					Accrual:    0,
					UploadedAt: orderUploadedAt,
				}

				pgxRowCreateOrder := pgxpool.NewRow(rowID).WithError(&pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation})
				mockPool.EXPECT().QueryRow(ctx, createOrder, userLogin, orderNumber).Return(pgxRowCreateOrder).Times(1)

				pgxRowGetOrder := pgxpool.NewRow(rowID, orderExpected.Login, orderExpected.Number, orderExpected.Status, orderExpected.Accrual, orderExpected.UploadedAt)
				mockPool.EXPECT().QueryRow(ctx, getOrder, orderNumber).Return(pgxRowGetOrder).Times(1)
			})

			It("returns the order and data conflict error", func() {
				result, err := repo.CreateOrder(ctx, &order)
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(repository.ErrConflict))
				Expect(*result).To(Equal(orderExpected))
			})
		})

		XWhen("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil order and an error", func() {

			})
		})
	})

	XContext("Calling GetListOfOrders method", func() {
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
				rowID = 1

				user = model.User{
					Login:    "user",
					Password: "password",
				}

				pgxRow := pgxpool.NewRow(rowID)
				mockPool.EXPECT().QueryRow(ctx, createBalance, userLogin).Return(pgxRow).Times(1)
			})

			It("returns nil error", func() {
				err = repo.CreateBalance(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		XWhen("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns an error", func() {

			})
		})
	})

	Context("Calling GetBalance method", func() {
		When("there is no error", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				userPassword = "password"
				var accrued float64 = 500
				var withdrawn float64 = 50

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				pgxRow := pgxpool.NewRow(rowID, userLogin, accrued, withdrawn)
				mockPool.EXPECT().QueryRow(ctx, getBalance, userLogin).Return(pgxRow).Times(1)
			})

			It("returns a balance and nil error", func() {
				result, err := repo.GetBalance(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Current).To(Equal(float64(450)))
				Expect(result.Withdrawn).To(Equal(float64(50)))
			})
		})

		XWhen("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil balance and an error", func() {

			})
		})
	})

	XContext("Calling WithdrawFromBalance method", func() {
		When("everything is right", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				orderNumber = "2377225624"

				var accrued float64 = 500
				var withdrawn float64 = 50
				var sum float64 = 100

				user = model.User{}

				mockPool.EXPECT().Begin(ctx).Times(1)
				pgxRowUpdate := pgxpool.NewRow(accrued, withdrawn)
				mockPool.EXPECT().QueryRow(ctx, updateBalanceWithdrawn, userLogin, withdrawn).Return(pgxRowUpdate).Times(1)
				pgxRowCreate := pgxpool.NewRow(rowID)
				mockPool.EXPECT().QueryRow(ctx, createWithdraw, userLogin, orderNumber, sum).Return(pgxRowCreate).Times(1)

				err = repo.WithdrawFromBalance(ctx, &user, orderNumber, sum)
			})

			It("returns nil error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("balance not enough to withdraw", func() {
			BeforeEach(func() {

			})

			It("returns negative balance error", func() {

			})
		})

		XWhen("something has gone wrong with the queries", func() {
			BeforeEach(func() {

			})

			It("returns an error", func() {

			})
		})
	})

	XContext("Calling GetListOfWithdrawals method", func() {
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

		XWhen("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil list of withdrawals and an error", func() {

			})
		})
	})

	XContext("Calling GetListOfOrdersToProcess method", func() {
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

		XWhen("something has gone wrong with the query", func() {
			BeforeEach(func() {

			})

			It("returns nil list of orders and an error", func() {

			})
		})
	})

	XContext("Calling UpdateBalanceAccrued method", func() {
		When("everything is right", func() {
			BeforeEach(func() {

			})

			It("returns nil error", func() {

			})
		})

		XWhen("something has gone wrong with the queries", func() {
			BeforeEach(func() {

			})

			It("returns an error", func() {

			})
		})
	})
})
