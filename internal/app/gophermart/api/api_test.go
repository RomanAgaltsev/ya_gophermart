package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/api"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/balance"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/order"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/user"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	balanceMocks "github.com/RomanAgaltsev/ya_gophermart/internal/mocks/balance"
	orderMocks "github.com/RomanAgaltsev/ya_gophermart/internal/mocks/order"
	userMocks "github.com/RomanAgaltsev/ya_gophermart/internal/mocks/user"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
	"github.com/RomanAgaltsev/ya_gophermart/internal/pkg/auth"

	"github.com/go-chi/jwtauth/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.uber.org/mock/gomock"
)

const (
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain; charset=utf-8"
)

var _ = Describe("Handler", func() {
	var (
		err                 error
		errSomethingStrange error

		cfg *config.Config

		server *ghttp.Server

		userService    user.Service
		userCtrl       *gomock.Controller
		userRepository *userMocks.MockRepository

		orderService    order.Service
		orderCtrl       *gomock.Controller
		orderRepository *orderMocks.MockRepository

		balanceService    balance.Service
		balanceCtrl       *gomock.Controller
		balanceRepository *balanceMocks.MockRepository

		handler *api.Handler

		endpoint string

		usr      *model.User
		usrBytes []byte

		ja     *jwtauth.JWTAuth
		cookie *http.Cookie

		expectOrders      model.Orders
		expectBalance     model.Balance
		expectWithdrawals model.Withdrawals

		withdrawal      model.Withdrawal
		withdrawalBytes []byte

		login       string
		secretKey   string
		tokenString string
		orderNumber string
	)

	BeforeEach(func() {
		errSomethingStrange = errors.New("something strange")

		cfg, err = config.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).ShouldNot(BeNil())

		server = ghttp.NewServer()

		// User service and repository
		userCtrl = gomock.NewController(GinkgoT())
		Expect(userCtrl).ShouldNot(BeNil())

		userRepository = userMocks.NewMockRepository(userCtrl)
		Expect(userRepository).ShouldNot(BeNil())

		userService, err = user.NewService(userRepository, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(userService).ShouldNot(BeNil())

		// Order service and repository
		orderCtrl = gomock.NewController(GinkgoT())
		Expect(orderCtrl).ShouldNot(BeNil())

		orderRepository = orderMocks.NewMockRepository(orderCtrl)
		Expect(orderRepository).ShouldNot(BeNil())

		orderService, err = order.NewService(orderRepository, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(orderService).ShouldNot(BeNil())

		// Balance service and repository
		balanceCtrl = gomock.NewController(GinkgoT())
		Expect(balanceCtrl).ShouldNot(BeNil())

		balanceRepository = balanceMocks.NewMockRepository(balanceCtrl)
		Expect(balanceRepository).ShouldNot(BeNil())

		balanceService, err = balance.NewService(balanceRepository, cfg, false)
		Expect(err).NotTo(HaveOccurred())
		Expect(balanceService).ShouldNot(BeNil())

		// Handler
		handler = api.NewHandler(cfg, userService, orderService, balanceService)
		Expect(handler).ShouldNot(BeNil())
	})

	AfterEach(func() {
		server.Close()
	})

	Context("Receiving request at the /api/user/register endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/register"
			server.AppendHandlers(handler.UserRegistrion)
		})

		When("the method is POST, content type is right and payload is right", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())

				userRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				balanceRepository.EXPECT().CreateBalance(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			})

			It("returns status 'OK' (200) and a cookie", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))

				cookie := resp.Header.Get("Set-Cookie")
				Expect(cookie).NotTo(BeEmpty())
			})
		})

		When("the method is POST, content type is right but payload is wrong", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Bad request' (400) and no cookie", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

				cookie := resp.Header.Get("Set-Cookie")
				Expect(cookie).To(BeEmpty())
			})
		})

		When("the method is POST, content type is wrong and payload is right", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Bad request' (400) and no cookie", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeText, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

				cookie := resp.Header.Get("Set-Cookie")
				Expect(cookie).To(BeEmpty())
			})
		})

		When("the method is POST, request is right but user already exists", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())

				userRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(repository.ErrConflict).Times(1)
			})

			It("returns status 'Conflict' (409) and no cookie", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusConflict))

				cookie := resp.Header.Get("Set-Cookie")
				Expect(cookie).To(BeEmpty())
			})
		})

		When("everything is right with the request, but something has gone wrong with the service", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())

				userRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("Receiving request at the /api/user/login endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/login"
			server.AppendHandlers(handler.UserLogin)
		})

		When("the method is POST, content type is right and payload is right", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				hash, err := auth.HashPassword("password")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(hash).ShouldNot(BeEmpty())

				expectUsr := &model.User{
					Login:    "user",
					Password: hash,
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())

				userRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(expectUsr, nil).Times(1)
			})

			It("returns status 'OK' (200) and a cookie", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))

				cookie := resp.Header.Get("Set-Cookie")
				Expect(cookie).NotTo(BeEmpty())
			})
		})

		When("the method is POST, content type is right but payload is wrong", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Bad request' (400) and no cookie", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

				cookie := resp.Header.Get("Set-Cookie")
				Expect(cookie).To(BeEmpty())
			})
		})

		When("the method is POST, content type is wrong and payload is right", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Bad request' (400) and no cookie", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeText, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))

				cookie := resp.Header.Get("Set-Cookie")
				Expect(cookie).To(BeEmpty())
			})
		})

		When("the method is POST, and login/password is wrong", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "wrong password",
				}

				hash, err := auth.HashPassword("password")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(hash).ShouldNot(BeEmpty())

				expectUsr := &model.User{
					Login:    "user",
					Password: hash,
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())

				userRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(expectUsr, nil).Times(1)
			})

			It("returns status 'Unauthorized' (401)", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusUnauthorized))
			})
		})

		When("everything is right with the request, but something has gone wrong with service", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				hash, err := auth.HashPassword("password")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(hash).ShouldNot(BeEmpty())

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())

				userRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				resp, err := http.Post(server.URL()+endpoint, ContentTypeJSON, bytes.NewReader(usrBytes))

				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("Receiving request at the /api/user/orders endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/register"
			server.RouteToHandler("POST", endpoint, handler.OrderNumberUpload)
			server.RouteToHandler("GET", endpoint, handler.OrderListRequest)

			secretKey = "secret"
			login = "user"

			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())

			_, tokenString, err = auth.NewJWTToken(ja, login)
			Expect(err).NotTo(HaveOccurred())
			Expect(tokenString).NotTo(BeEmpty())

			cookie = auth.NewCookieWithDefaults(tokenString)
		})

		// POST
		When("the method is POST, everything is right but the order has been already created by this user", func() {
			BeforeEach(func() {
				orderNumber = "12345678903"

				expectOrder := &model.Order{
					Login:      login,
					Number:     orderNumber,
					Status:     "NEW",
					Accrual:    0,
					UploadedAt: time.Now(),
				}

				orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(expectOrder, repository.ErrConflict).Times(1)
			})

			It("returns status 'OK' (200)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader([]byte(orderNumber)))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeText)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))
			})
		})

		When("the method is POST, everything is right and order doesn`t exist", func() {
			BeforeEach(func() {
				orderNumber = "12345678903"

				expectOrder := &model.Order{
					Login:      login,
					Number:     orderNumber,
					Status:     "NEW",
					Accrual:    0,
					UploadedAt: time.Now(),
				}

				orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(expectOrder, nil).Times(1)
			})

			It("returns status 'Accepted' (202)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader([]byte(orderNumber)))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeText)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusAccepted))
			})
		})

		When("the method is POST and order number is empty", func() {
			BeforeEach(func() {
				orderNumber = ""
			})

			It("returns status 'Bad request' (400)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader([]byte(orderNumber)))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeText)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusBadRequest))
			})
		})

		When("the method is POST, everything is right but the order has been already created by another user", func() {
			BeforeEach(func() {
				orderNumber = "12345678903"

				expectOrder := &model.Order{
					Login:      "another user",
					Number:     orderNumber,
					Status:     "NEW",
					Accrual:    0,
					UploadedAt: time.Now(),
				}

				orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(expectOrder, repository.ErrConflict).Times(1)
			})

			It("returns status 'Conflict' (409)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader([]byte(orderNumber)))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeText)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusConflict))
			})
		})

		When("the method is POST and the order number is invalid", func() {
			BeforeEach(func() {
				orderNumber = "order #123456"
			})

			It("returns status 'Unprocessable entity' (422)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader([]byte(orderNumber)))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeText)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusUnprocessableEntity))
			})
		})

		When("the method is POST, everything is right with the request, but something has gone wrong with service", func() {
			BeforeEach(func() {
				orderNumber = "12345678903"

				orderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil, errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader([]byte(orderNumber)))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeText)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})

		// GET
		When("the method is GET and there are orders to return", func() {
			BeforeEach(func() {
				expectOrders = []*model.Order{
					{
						Login:      login,
						Number:     orderNumber,
						Status:     "NEW",
						Accrual:    0,
						UploadedAt: time.Now(),
					},
				}

				orderRepository.EXPECT().GetListOfOrders(gomock.Any(), gomock.Any()).Return(expectOrders, nil).Times(1)
			})

			It("returns status 'OK' (200) and a list of orders in JSON", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var orders model.Orders
				err = json.NewDecoder(response.Body).Decode(&orders)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(orders).ShouldNot(BeEmpty())
				Expect(orders).Should(HaveLen(len(expectOrders)))
			})
		})

		When("the method is GET and there are no orders to return", func() {
			BeforeEach(func() {
				expectOrders = []*model.Order{}

				orderRepository.EXPECT().GetListOfOrders(gomock.Any(), gomock.Any()).Return(expectOrders, nil).Times(1)
			})

			It("returns status 'No content' (204) and response body is empty", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusNoContent))
			})
		})

		When("the method is GET, but something has gone wrong with the service", func() {
			BeforeEach(func() {
				orderRepository.EXPECT().GetListOfOrders(gomock.Any(), gomock.Any()).Return(nil, errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("Receiving request at the /api/user/balance endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/balance"
			server.AppendHandlers(handler.UserBalanceRequest)

			secretKey = "secret"
			login = "user"

			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())

			_, tokenString, err = auth.NewJWTToken(ja, login)
			Expect(err).NotTo(HaveOccurred())
			Expect(tokenString).NotTo(BeEmpty())

			cookie = auth.NewCookieWithDefaults(tokenString)
		})

		When("the method is GET and everything is right", func() {
			BeforeEach(func() {
				expectBalance = model.Balance{
					Current:   500,
					Withdrawn: 42,
				}

				balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(&expectBalance, nil).Times(1)
			})

			It("returns status 'OK' (200) and a balance structure in JSON", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var blnc model.Balance
				err = json.NewDecoder(response.Body).Decode(&blnc)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(blnc).To(Equal(expectBalance))
			})
		})

		When("the method is GET, but something has gone wrong with the service", func() {
			BeforeEach(func() {
				balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(nil, errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("Receiving request at the /api/user/balance/withdraw endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/balance/withdraw"
			server.AppendHandlers(handler.WithdrawRequest)

			secretKey = "secret"
			login = "user"

			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())

			_, tokenString, err = auth.NewJWTToken(ja, login)
			Expect(err).NotTo(HaveOccurred())
			Expect(tokenString).NotTo(BeEmpty())

			cookie = auth.NewCookieWithDefaults(tokenString)
		})

		When("the method is POST and balance is enough to withdraw", func() {
			BeforeEach(func() {
				withdrawal = model.Withdrawal{
					Login:       login,
					OrderNumber: "2377225624",
					Sum:         751,
					ProcessedAt: time.Now(),
				}

				withdrawalBytes, err = json.Marshal(withdrawal)
				Expect(err).ShouldNot(HaveOccurred())

				balanceRepository.EXPECT().WithdrawFromBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
			})

			It("returns status 'OK' (200)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(withdrawalBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))
			})
		})

		When("the method is POST and balance is not enough to withdraw", func() {
			BeforeEach(func() {
				withdrawal = model.Withdrawal{
					OrderNumber: "2377225624",
					Sum:         751,
				}

				withdrawalBytes, err = json.Marshal(withdrawal)
				Expect(err).ShouldNot(HaveOccurred())

				balanceRepository.EXPECT().WithdrawFromBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(repository.ErrNegativeBalance).Times(1)
			})

			It("returns status 'Payment required' (402)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(withdrawalBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusPaymentRequired))
			})
		})

		When("the method is POST and the order number is invalid", func() {
			BeforeEach(func() {
				withdrawal = model.Withdrawal{
					OrderNumber: "Order #12345",
					Sum:         751,
				}

				withdrawalBytes, err = json.Marshal(withdrawal)
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("returns status 'Unprocessable entity' (422)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(withdrawalBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusUnprocessableEntity))
			})
		})

		When("the method is POST, everything is right with the request, but something has gone wrong with service", func() {
			BeforeEach(func() {
				withdrawal = model.Withdrawal{
					OrderNumber: "2377225624",
					Sum:         751,
				}

				withdrawalBytes, err = json.Marshal(withdrawal)
				Expect(err).ShouldNot(HaveOccurred())

				balanceRepository.EXPECT().WithdrawFromBalance(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				request, err := http.NewRequest(http.MethodPost, server.URL()+endpoint, bytes.NewReader(withdrawalBytes))
				Expect(err).ShouldNot(HaveOccurred())

				request.Header.Set("Content-Type", ContentTypeJSON)
				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})

	Context("Receiving request at the /api/user/withdrawals endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/withdrawals"
			server.AppendHandlers(handler.WithdrawalsInformationRequest)

			secretKey = "secret"
			login = "user"

			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())

			_, tokenString, err = auth.NewJWTToken(ja, login)
			Expect(err).NotTo(HaveOccurred())
			Expect(tokenString).NotTo(BeEmpty())

			cookie = auth.NewCookieWithDefaults(tokenString)
		})

		When("the method is GET and there are withdrawals to return", func() {
			BeforeEach(func() {
				expectWithdrawals = model.Withdrawals{
					{
						OrderNumber: "2377225624",
						Sum:         500,
						ProcessedAt: time.Now(),
					},
				}

				balanceRepository.EXPECT().GetListOfWithdrawals(gomock.Any(), gomock.Any()).Return(expectWithdrawals, nil).Times(1)
			})

			It("returns status 'OK' (200) and a list of withdrawals in JSON", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusOK))

				var withdrawals model.Withdrawals
				err = json.NewDecoder(response.Body).Decode(&withdrawals)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(expectWithdrawals).ShouldNot(BeEmpty())
				Expect(withdrawals).Should(HaveLen(len(expectWithdrawals)))
			})
		})

		When("the method is GET and there are no withdrawals to return", func() {
			BeforeEach(func() {
				expectWithdrawals = model.Withdrawals{}

				balanceRepository.EXPECT().GetListOfWithdrawals(gomock.Any(), gomock.Any()).Return(expectWithdrawals, nil).Times(1)
			})

			It("returns status 'No content' (204) and response body is empty", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusNoContent))
			})
		})

		When("the method is GET, but something has gone wrong with the service", func() {
			BeforeEach(func() {
				balanceRepository.EXPECT().GetListOfWithdrawals(gomock.Any(), gomock.Any()).Return(nil, errSomethingStrange).Times(1)
			})

			It("returns status 'Internal server error' (500)", func() {
				request, err := http.NewRequest(http.MethodGet, server.URL()+endpoint, nil)
				Expect(err).ShouldNot(HaveOccurred())

				request.AddCookie(cookie)

				response, err := http.DefaultClient.Do(request)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(response.StatusCode).Should(Equal(http.StatusInternalServerError))
			})
		})
	})
})
