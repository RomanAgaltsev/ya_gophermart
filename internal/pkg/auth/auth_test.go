package auth_test

import (
	"github.com/RomanAgaltsev/ya_gophermart/internal/pkg/auth"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	Describe("Creating new JWT token", func() {
		var (
			ja        *jwtauth.JWTAuth
			secretKey string
			login     string
		)

		JustBeforeEach(func() {
			ja = auth.NewAuth(secretKey)
			Expect(ja).ShouldNot(BeNil())
		})

		Context("When the secret key and login are defined and correct", func() {
			BeforeEach(func() {
				secretKey = "secret"
				login = "user"
			})

			It("can create new valid JWT token", func() {
				token, tokenString, err := auth.NewJWTToken(ja, login)
				Expect(err).NotTo(HaveOccurred())
				Expect(tokenString).NotTo(BeEmpty())

				err = jwt.Validate(token, ja.ValidateOptions()...)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("When the secret key and login are defined but not correct", func() {
			BeforeEach(func() {
				secretKey = ""
				login = ""
			})

			It("cannot create new JWT token", func() {
				token, tokenString, err := auth.NewJWTToken(ja, login)
				Expect(err).To(HaveOccurred())
				Expect(tokenString).To(BeEmpty())

				Expect(func() {
					err = jwt.Validate(token, ja.ValidateOptions()...)
				}).To(Panic())
			})
		})
	})

	Describe("Creating new cookie", func() {

	})

	Describe("Hashing password", func() {

	})

	Describe("Extracting user login from HTTP request", func() {

	})
})
