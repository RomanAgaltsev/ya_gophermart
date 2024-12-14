package auth_test

import (
	"github.com/RomanAgaltsev/ya_gophermart/internal/model"
	"net/http"

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

		Context("When the secret key and login are undefined", func() {
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
		var (
			value  string
			cookie *http.Cookie
		)

		Context("When the cookie value is defined", func() {
			BeforeEach(func() {
				value = "value"
			})

			It("can create new cookie with a given value", func() {
				cookie = auth.NewCookieWithDefaults(value)
				Expect(cookie.Name).To(Equal(auth.DefaultCookieName))
				Expect(cookie.Value).To(Equal(value))
				Expect(cookie.Path).To(Equal(auth.DefaultCookiePath))
				Expect(cookie.MaxAge).To(Equal(auth.DefaultCookieMaxAge))
				Expect(cookie.SameSite).To(Equal(http.SameSiteDefaultMode))
			})
		})

		Context("When the cookie value is undefined", func() {
			BeforeEach(func() {
				value = ""
			})

			It("can create new cookie without any value", func() {
				cookie = auth.NewCookieWithDefaults(value)
				Expect(cookie.Name).To(Equal(auth.DefaultCookieName))
				Expect(cookie.Value).To(BeEmpty())
				Expect(cookie.Path).To(Equal(auth.DefaultCookiePath))
				Expect(cookie.MaxAge).To(Equal(auth.DefaultCookieMaxAge))
				Expect(cookie.SameSite).To(Equal(http.SameSiteDefaultMode))
			})
		})
	})

	Describe("Hashing password", func() {
		var (
			password string
			hash     string
			err      error
		)

		Context("When the password is defined", func() {
			BeforeEach(func() {
				password = "password"
			})

			It("can hash the password", func() {
				hash, err = auth.HashPassword(password)
				Expect(err).NotTo(HaveOccurred())
				Expect(hash).NotTo(BeEmpty())
			})
			It("can check the password hash", func() {
				result := auth.CheckPasswordHash(password, hash)
				Expect(result).To(BeTrue())
			})
		})

		Context("When the password is undefined", func() {
			BeforeEach(func() {
				password = ""
			})

			It("can hash the password", func() {
				hash, err = auth.HashPassword(password)
				Expect(err).NotTo(HaveOccurred())
				Expect(hash).NotTo(BeEmpty())
			})
			It("can check the password hash", func() {
				result := auth.CheckPasswordHash(password, hash)
				Expect(result).To(BeTrue())
			})
		})

		Context("When the password is longer than 72 bytes", func() {
			BeforeEach(func() {
				password = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFJHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyzABCDEFJHIJKLMNOPQRSTUVWXYZ"

			})

			It("cannot hash the password", func() {
				hash, err = auth.HashPassword(password)
				Expect(err).To(HaveOccurred())
				Expect(hash).To(BeEmpty())
			})
		})
	})

	Describe("Extracting user login from HTTP request", func() {
		const secretKeyEnc = "secret"

		var (
			ja          *jwtauth.JWTAuth
			cookie      *http.Cookie
			request     *http.Request
			user        *model.User
			login       string
			secretKey   string
			tokenString string
			err         error
		)

		JustBeforeEach(func() {
			ja = auth.NewAuth(secretKeyEnc)
			Expect(ja).ShouldNot(BeNil())

			_, tokenString, err = auth.NewJWTToken(ja, login)
			Expect(err).NotTo(HaveOccurred())
			Expect(tokenString).NotTo(BeEmpty())

			cookie = auth.NewCookieWithDefaults(tokenString)

			request, err = http.NewRequest("", "", nil)
			Expect(err).NotTo(HaveOccurred())

			request.AddCookie(cookie)
		})

		Context("When the secret key is defined and right", func() {
			BeforeEach(func() {
				secretKey = secretKeyEnc
				login = "user"
			})

			It("can get user from request", func() {
				user, err = auth.UserFromRequest(request, secretKey)
				Expect(err).NotTo(HaveOccurred())
				Expect(user.Login).To(Equal(login))
			})
		})

		Context("When the secret key is defined and wrong", func() {
			BeforeEach(func() {
				secretKey = "wrong key"
				login = "user"
			})

			It("cannot get user from request", func() {
				user, err = auth.UserFromRequest(request, secretKey)
				Expect(err).To(HaveOccurred())
				Expect(user).To(BeNil())
			})
		})

		Context("When the secret key is undefined", func() {
			BeforeEach(func() {
				secretKey = ""
				login = "user"
			})

			It("cannot get user from request", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user).To(BeNil())
			})
		})
	})
})
