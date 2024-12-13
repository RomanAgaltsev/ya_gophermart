package api_test

/*
import (
	"bytes"
	"encoding/json"
	"net/http"

	//"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/api"
	//"github.com/RomanAgaltsev/ya_gophermart/internal/config"
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Handler", func() {
	var (
		server *ghttp.Server

		//	cfg, _  = config.Get()
		//handler = api.NewHandler(cfg)

		testUser = model.User{
			Login:    "user",
			Password: "password",
		}

		userBytes, _ = json.Marshal(testUser)
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
	})
	AfterEach(func() {
		server.Close()
	})

	Context("When POST request is sent to the /api/user/register path", func() {
		//		BeforeEach(func() {
		//			server.AppendHandlers(handler.UserRegistrion)
		//		})

		It("Returns status OK (200) and a cookie", func() {
			resp, err := http.Post(server.URL()+"/api/user/register", "application/json", bytes.NewReader(userBytes))

			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			//resp.Cookies()
		})
	})

})
*/
