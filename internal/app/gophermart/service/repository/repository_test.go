package repository_test

import (
	"context"
	"errors"
	//"github.com/jackc/pgx/v5"
	"time"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pashagolub/pgxmock/v4"
)

var _ = Describe("Repository", func() {
	var (
		err                 error
		errSomethingStrange error

		ctx context.Context

		mockPool pgxmock.PgxPoolIface
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

		mockPool, err = pgxmock.NewPool()
		Expect(err).ShouldNot(HaveOccurred())

		repo, err = repository.New(mockPool)
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		mockPool.Close()
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

				rs := pgxmock.NewRows([]string{"id"}).
					AddRow(rowID)
				mockPool.ExpectQuery("INSERT .+ VALUES .+").
					WithArgs(userLogin, userPassword).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
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

				rs := pgxmock.NewRows([]string{"id"}).
					AddRow(rowID).
					RowError(int(rowID), &pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation})
				mockPool.ExpectQuery("INSERT .+ VALUES .+").
					WithArgs(userLogin, userPassword).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
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

				rs := pgxmock.NewRows([]string{"id", "login", "password", "createdat"}).
					AddRow(rowID, userLogin, userPassword, userCreatedAt)
				mockPool.ExpectQuery("SELECT .+ FROM users WHERE .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
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

				rs := pgxmock.NewRows([]string{"id", "login", "password", "createdat"}).
					AddRow(rowID, userLogin, userPassword, userCreatedAt)
				mockPool.ExpectQuery("SELECT .+ FROM users WHERE .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns a user and nil error", func() {
				result, err := repo.GetUser(ctx, userLogin)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(*result).To(Equal(userExpected))
			})
		})

		// Pending - because of exponential backoff it lasts about 10 minutes
		// Remove 'X' to run
		XWhen("something has gone wrong with the query", func() {
			BeforeEach(func() {
				rowID = 0
				userLogin = ""
				userPassword = ""
				userCreatedAt = time.Now()

				rs := pgxmock.NewRows([]string{"id", "login", "password", "createdat"}).
					RowError(1, errSomethingStrange)
				mockPool.ExpectQuery("SELECT .+ FROM users WHERE .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
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

				rs := pgxmock.NewRows([]string{"id"}).
					AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO orders .+ VALUES .+").
					WithArgs(userLogin, orderNumber).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
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

				rsCreate := pgxmock.NewRows([]string{"id"}).
					AddRow(rowID).
					RowError(int(rowID), &pgconn.PgError{Code: pgerrcode.IntegrityConstraintViolation})
				mockPool.ExpectQuery("INSERT INTO orders .+ VALUES .+").
					WithArgs(userLogin, orderNumber).
					WillReturnRows(rsCreate).
					Times(1)

				rsGet := pgxmock.NewRows([]string{"id", "login", "ordernumber", "status", "accrual", "uploadedat"}).
					AddRow(rowID, orderExpected.Login, orderExpected.Number, orderExpected.Status, orderExpected.Accrual, orderExpected.UploadedAt)
				mockPool.ExpectQuery("SELECT .+ FROM orders .+").
					WithArgs(orderNumber).
					WillReturnRows(rsGet).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns the order and data conflict error", func() {
				result, err := repo.CreateOrder(ctx, &order)
				Expect(err).Should(HaveOccurred())
				Expect(err).To(Equal(repository.ErrConflict))
				Expect(*result).To(Equal(orderExpected))
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
	})

	Context("Calling CreateBalance method", func() {
		When("the balance is successfully created", func() {
			BeforeEach(func() {
				rowID = 1

				user = model.User{
					Login:    "user",
					Password: "password",
				}

				rs := pgxmock.NewRows([]string{"id"}).
					AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO balance .+ VALUES .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns nil error", func() {
				err = repo.CreateBalance(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
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

				rs := pgxmock.NewRows([]string{"id", "login", "accrued", "withdrawn"}).
					AddRow(rowID, userLogin, accrued, withdrawn)
				mockPool.ExpectQuery("SELECT .+ FROM balance .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})

			It("returns a balance and nil error", func() {
				result, err := repo.GetBalance(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Current).To(Equal(float64(450)))
				Expect(result.Withdrawn).To(Equal(float64(50)))
			})
		})
	})

	Context("Calling WithdrawFromBalance method", func() {
		When("everything is right", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				orderNumber = "2377225624"

				//				var accrued float64 = 500
				//				var withdrawn float64 = 50
				var sum float64 = 100
				accrued := 500
				withdrawn := 50
				//sum := 100

				user = model.User{}

				//				mockPool.EXPECT().Begin(ctx).Times(1)
				//				pgxRowUpdate := pgxpool.NewRow(accrued, withdrawn)
				//				mockPool.EXPECT().QueryRow(ctx, updateBalanceWithdrawn, userLogin, withdrawn).Return(pgxRowUpdate).Times(1)
				//				pgxRowCreate := pgxpool.NewRow(rowID)
				//				mockPool.EXPECT().QueryRow(ctx, createWithdraw, userLogin, orderNumber, sum).Return(pgxRowCreate).Times(1)

				//mockPool.ExpectBeginTx(pgx.TxOptions{AccessMode: pgx.ReadOnly})
				//mockPool.ExpectBeginTx(pgx.TxOptions{})
				mockPool.ExpectBegin()

				rsUpdate := pgxmock.NewRows([]string{"accrued", "withdrawn"}).
					AddRow(accrued, withdrawn)
				mockPool.ExpectQuery("UPDATE balance .+ SET .+").
					WithArgs(userLogin, withdrawn).
					WillReturnRows(rsUpdate).
					Times(1)

				rsCreate := pgxmock.NewRows([]string{"id"}).
					AddRow(rowID)
				mockPool.ExpectQuery("INSERT INTO withdrawals .+ VALUES .+").
					WithArgs(userLogin, orderNumber, sum).
					WillReturnRows(rsCreate).
					Times(1)

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()

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
