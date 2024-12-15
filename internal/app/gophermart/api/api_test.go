package api_test

import (
	"bytes"
	"errors"
	"github.com/RomanAgaltsev/ya_gophermart/internal/pkg/auth"

	//"context"
	"encoding/json"
	"net/http"

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
		err error

		//ctx context.Context
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
	)

	BeforeEach(func() {
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

		balanceService, err = balance.NewService(balanceRepository, cfg)
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
		/*
			When("the method is GET", func() {
				It("returns status 'Method not allowed' (405)", func() {
					resp, err := http.Get(server.URL() + endpoint)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.StatusCode).Should(Equal(http.StatusMethodNotAllowed))
				})
			})
		*/

		When("everything is right with the request, but something has gone wrong with the service", func() {
			BeforeEach(func() {
				usr = &model.User{
					Login:    "user",
					Password: "password",
				}

				usrBytes, err = json.Marshal(usr)
				Expect(err).ShouldNot(HaveOccurred())

				userRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("a strange mistake")).Times(1)
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

		//		When("the method is GET and no matter what content type and payload", func() {
		//			It("returns status 'Method not allowed' (405)", func() {
		//
		//			})
		//		})

		When("everything with the request is right, but something has gone wrong with service", func() {
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

				userRepository.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("a strange mistake")).Times(1)
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

			//server.AppendHandlers(handler.OrderNumberUpload)
			//server.AppendHandlers(handler.OrderListRequest)
		})
	})

	Context("Receiving request at the /api/user/balance endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/balance"
			server.AppendHandlers(handler.UserBalanceRequest)
		})
	})

	Context("Receiving request at the /api/user/balance/withdraw endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/balance/withdraw"
			server.AppendHandlers(handler.WithdrawRequest)
		})
	})

	Context("Receiving request at the /api/user/withdrawals endpoint", func() {
		BeforeEach(func() {
			endpoint = "/api/user/withdrawals"
			server.AppendHandlers(handler.WithdrawalsInformationRequest)
		})
	})

})
