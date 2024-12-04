package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type UserLogin string

const (
	// DefaultCookieName содержит имя куки по умолчанию.
	DefaultCookieName = "jwt"

	// DefaultCookiePath содержит путь куки по умолчанию.
	DefaultCookiePath = "/"

	// DefaultCookieMaxAge содержит возраст куки по умолчанию.
	DefaultCookieMaxAge = 3600

	// UserLoginClaimName содержит имя ключа логина пользователя в контексте.
	UserLoginClaimName UserLogin = "login"
)

// NewJWTToken создает новый JWT токен.
func NewJWTToken(ja *jwtauth.JWTAuth, login string) (token jwt.Token, tokenString string, err error) {
	return ja.Encode(map[string]interface{}{string(UserLoginClaimName): login})
}

// NewCookieWithDefaults создает новую куку со значениями по умолчанию и переданным в параметре значением.
func NewCookieWithDefaults(value string) *http.Cookie {
	return &http.Cookie{
		Name:     DefaultCookieName,
		Value:    value,
		Path:     DefaultCookiePath,
		MaxAge:   DefaultCookieMaxAge,
		SameSite: http.SameSiteDefaultMode,
	}
}
