package repository_test

import (
	"context"
	"errors"
	"github.com/RomanAgaltsev/ya_gophermart/internal/database/queries"
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

		// Withdrawal
		withdrawalProcessedAt time.Time

		// Accrual
		orderAccrual model.OrderAccrual
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
					Status:     queries.OrderStatusNEW,
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

	Context("Calling GetListOfOrders method", func() {
		When("orders exist", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				userPassword = "password"
				orderStatus := queries.OrderStatusNEW
				var accrual float64 = 100
				orderUploadedAt = time.Now()

				rs := pgxmock.NewRows([]string{"id", "login", "ordernumber", "status", "accrual", "uploadedat"}).
					AddRow(rowID, userLogin, orderNumber, orderStatus, accrual, orderUploadedAt)
				mockPool.ExpectQuery("SELECT .+ FROM orders .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns a non-empty list of orders and nil error", func() {
				result, err := repo.GetListOfOrders(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(Equal(1))
			})
		})

		When("orders don't exist", func() {
			BeforeEach(func() {
				userLogin = "user"

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				rs := pgxmock.NewRows([]string{"id", "login", "ordernumber", "status", "accrual", "uploadedat"})
				mockPool.ExpectQuery("SELECT .+ FROM orders .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)

			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an empty list of orders and nil error", func() {
				result, err := repo.GetListOfOrders(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(Equal(0))
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
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
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
				userPassword = "password"
				orderNumber = "2377225624"

				var accrued float64 = 500
				var withdrawn float64 = 50
				var sum float64 = 100

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				mockPool.ExpectBegin()

				rsUpdate := pgxmock.NewRows([]string{"accrued", "withdrawn"}).
					AddRow(accrued, withdrawn)
				mockPool.ExpectQuery("UPDATE balance SET .+").
					WithArgs(userLogin, sum).
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
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns nil error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("balance not enough to withdraw", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				userPassword = "password"
				orderNumber = "2377225624"

				var accrued float64 = 100
				var withdrawn float64 = 150
				var sum float64 = 100

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				mockPool.ExpectBegin()

				rsUpdate := pgxmock.NewRows([]string{"accrued", "withdrawn"}).
					AddRow(accrued, withdrawn)
				mockPool.ExpectQuery("UPDATE balance SET .+").
					WithArgs(userLogin, sum).
					WillReturnRows(rsUpdate).
					Times(1)

				mockPool.ExpectRollback()

				err = repo.WithdrawFromBalance(ctx, &user, orderNumber, sum)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns negative balance error", func() {
				Expect(err).To(Equal(repository.ErrNegativeBalance))
			})
		})
	})

	Context("Calling GetListOfWithdrawals method", func() {
		When("withdrawals exist", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				userPassword = "password"
				var sum float64 = 50

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				rs := pgxmock.NewRows([]string{"id", "login", "ordernumber", "sum", "processedat"}).
					AddRow(rowID, userLogin, orderNumber, sum, withdrawalProcessedAt)
				mockPool.ExpectQuery("SELECT .+ FROM withdrawals .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns a non-empty list of withdrawals and nil error", func() {
				result, err := repo.GetListOfWithdrawals(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(Equal(1))
			})
		})

		When("withdrawals don't exist", func() {
			BeforeEach(func() {
				userLogin = "user"

				user = model.User{
					Login:    userLogin,
					Password: userPassword,
				}

				rs := pgxmock.NewRows([]string{"id", "login", "ordernumber", "sum", "processedat"})
				mockPool.ExpectQuery("SELECT .+ FROM withdrawals .+").
					WithArgs(userLogin).
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an empty list of withdrawals and nil error", func() {
				result, err := repo.GetListOfWithdrawals(ctx, &user)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(Equal(0))
			})
		})
	})

	Context("Calling GetListOfOrdersToProcess method", func() {
		When("orders to process exist", func() {
			BeforeEach(func() {
				rowID = 1
				userLogin = "user"
				orderNumber = "12345678903"
				orderStatus := queries.OrderStatusNEW
				var accrual float64 = 100
				orderUploadedAt = time.Now()

				rs := pgxmock.NewRows([]string{"id", "login", "ordernumber", "status", "accrual", "uploadedat"}).
					AddRow(rowID, userLogin, orderNumber, orderStatus, accrual, orderUploadedAt)
				mockPool.ExpectQuery("SELECT .+ FROM orders .+").
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns a non-empty list of orders and nil error", func() {
				result, err := repo.GetListOfOrdersToProcess(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(Equal(1))
			})
		})

		When("orders to process don't exist", func() {
			BeforeEach(func() {

				rs := pgxmock.NewRows([]string{"id", "login", "ordernumber", "status", "accrual", "uploadedat"})
				mockPool.ExpectQuery("SELECT .+ FROM orders .+").
					WillReturnRows(rs).
					Times(1)
			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns an empty list of orders and nil error", func() {
				result, err := repo.GetListOfOrdersToProcess(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(len(result)).To(Equal(0))
			})
		})
	})

	Context("Calling UpdateBalanceAccrued method", func() {
		When("everything is right", func() {
			BeforeEach(func() {
				var accrued float64 = 100
				var withdrawn float64 = 0

				userLogin = "user"
				orderNumber = "12345678903"
				orderStatus := queries.OrderStatusNEW

				order = model.Order{
					Login:      userLogin,
					Number:     orderNumber,
					Status:     orderStatus,
					Accrual:    0,
					UploadedAt: time.Now(),
				}
				orderAccrual = model.OrderAccrual{
					OrderNumber: orderNumber,
					Status:      queries.OrderStatusNEW,
					Accrual:     accrued,
				}

				mockPool.ExpectBegin()

				rs := pgxmock.NewRows([]string{"accrued", "withdrawn"}).
					AddRow(accrued, withdrawn)
				mockPool.ExpectQuery("UPDATE balance SET .+").
					WithArgs(userLogin, accrued).
					WillReturnRows(rs).
					Times(1)

				resultOrders := pgxmock.NewResult("UPDATE", 1)
				mockPool.ExpectExec("UPDATE orders SET .+").
					WithArgs(orderNumber, orderStatus, accrued).
					WillReturnResult(resultOrders).
					Times(1)

				mockPool.ExpectCommit()
				mockPool.ExpectRollback()

			})
			AfterEach(func() {
				err = mockPool.ExpectationsWereMet()
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns nil error", func() {
				err = repo.UpdateBalanceAccrued(ctx, &order, &orderAccrual)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
