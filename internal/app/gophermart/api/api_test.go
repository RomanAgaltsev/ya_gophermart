package api_test

import (
	"bytes"
	"context"

	//"encoding/json"

	"net/http"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/api"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/balance"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/order"
	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/user"
	"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	balanceMocks "github.com/RomanAgaltsev/ya_gophermart/internal/mocks/balance"
	orderMocks "github.com/RomanAgaltsev/ya_gophermart/internal/mocks/order"
	userMocks "github.com/RomanAgaltsev/ya_gophermart/internal/mocks/user"
	//"github.com/RomanAgaltsev/ya_gophermart/internal/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Handler", func() {
	var (
		err error

		ctx context.Context
		cfg *config.Config

		server *ghttp.Server

		userService    user.Service
		orderService   order.Service
		balanceService balance.Service

		userCtrl       *gomock.Controller
		userRepository *userMocks.MockRepository

		orderCtrl       *gomock.Controller
		orderRepository *orderMocks.MockRepository

		balanceCtrl       *gomock.Controller
		balanceRepository *balanceMocks.MockRepository

		handler *api.Handler
	)

	BeforeEach(func() {
		ctx = context.Background()

		cfg, err = config.Get()
		Expect(err).NotTo(HaveOccurred())
		Expect(err).ShouldNot(BeNil())

		server = ghttp.NewServer()

		// User
		userCtrl = gomock.NewController(GinkgoT())
		Expect(userCtrl).ShouldNot(BeNil())

		userRepository = userMocks.NewMockRepository(userCtrl)
		Expect(userRepository).ShouldNot(BeNil())

		userService, err = user.NewService(userRepository, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(userService).ShouldNot(BeNil())

		// Order
		orderCtrl = gomock.NewController(GinkgoT())
		Expect(orderCtrl).ShouldNot(BeNil())

		orderRepository = orderMocks.NewMockRepository(orderCtrl)
		Expect(orderRepository).ShouldNot(BeNil())

		orderService, err = order.NewService(orderRepository, cfg)
		Expect(err).NotTo(HaveOccurred())
		Expect(orderService).ShouldNot(BeNil())

		// Balance
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
			server.AppendHandlers(handler.UserRegistrion)
		})

		When("the method is POST, content type is right and payload is right", func() {
			It("returns status 'OK' (200) and a cookie", func() {

			})
		})

		When("the method is POST, content type is right but payload is wrong", func() {
			It("returns status 'Bad request' (400) and no cookie", func() {

			})
		})

		When("the method is POST, content type is wrong and payload is wrong", func() {
			It("returns status 'Bad request' (400) and no cookie", func() {

			})
		})

		When("the method is POST, request is right but user already exists", func() {
			It("returns status 'Conflict' (409) and no cookie", func() {

			})
		})

		When("the method is GET and no matter what content type and payload", func() {
			It("returns status 'Method not allowed' (405)", func() {

			})
		})

		When("everything with the request is right, but something has gone wrong with the service", func() {
			It("returns status 'Internal server error' (500)", func() {

			})
		})

		//		When("", func() {
		//			It("Returns status OK (200) and a cookie", func() {
		//				//resp, err := http.Post(server.URL()+"/api/user/register", "application/json", bytes.NewReader(userBytes))
		//				resp, err := http.Post(server.URL()+"/api/user/register", "application/json", nil)
		//
		//				Expect(err).ShouldNot(HaveOccurred())
		//				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		//
		//				//resp.Cookies()
		//			})
		//		})
	})

	Context("Receiving request at the /api/user/login endpoint", func() {
		BeforeEach(func() {
			server.AppendHandlers(handler.UserLogin)
		})

		When("the method is POST, content type is right and payload is right", func() {
			It("returns status 'OK' (200) and a cookie", func() {

			})
		})

		When("the method is POST, content type is right but payload is wrong", func() {
			It("returns status 'Bad request' (400) and no cookie", func() {

			})
		})

		When("the method is POST, content type is wrong and payload is wrong", func() {
			It("returns status 'Bad request' (400) and no cookie", func() {

			})
		})

		// !!Middleware
		When("the method is POST, but login/password is wrong", func() {
			It("returns status 'Unauthorized' (401)", func() {

			})
		})

		When("the method is GET and no matter what content type and payload", func() {
			It("returns status 'Method not allowed' (405)", func() {

			})
		})

		When("everything with the request is right, but something has gone wrong with service", func() {
			It("returns status 'Internal server error' (500)", func() {

			})
		})
	})

	Context("Receiving request at the /api/user/orders endpoint", func() {
		BeforeEach(func() {
			server.AppendHandlers(handler.OrderNumberUpload)
			server.AppendHandlers(handler.OrderListRequest)
		})
	})

	Context("Receiving request at the /api/user/balance endpoint", func() {
		BeforeEach(func() {
			server.AppendHandlers(handler.UserBalanceRequest)
		})
	})

	Context("Receiving request at the /api/user/balance/withdraw endpoint", func() {
		BeforeEach(func() {
			server.AppendHandlers(handler.WithdrawRequest)
		})
	})

	Context("Receiving request at the /api/user/withdrawals endpoint", func() {
		BeforeEach(func() {
			server.AppendHandlers(handler.WithdrawalsInformationRequest)
		})
	})

})
