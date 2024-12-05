package auth

import (
    "net/http"

    "github.com/go-chi/jwtauth/v5"
    "github.com/lestrrat-go/jwx/v2/jwt"
    "golang.org/x/crypto/bcrypt"
)

type UserLogin string

const (
    // JWTSignAlgorithm contains JWT signing algorithm
    JWTSignAlgorithm = "HS256"

    // DefaultCookieName contains default cookie name.
    DefaultCookieName = "jwt"

    // DefaultCookiePath contains default cookie path.
    DefaultCookiePath = "/"

    // DefaultCookieMaxAge contains default cookie max age.
    DefaultCookieMaxAge = 3600

    // UserLoginClaimName contains key name of user login in a context.
    UserLoginClaimName UserLogin = "login"
)

// NewAuth returns new JWTAuth.
func NewAuth(secretKey string) *jwtauth.JWTAuth {
    return jwtauth.New(JWTSignAlgorithm, []byte(secretKey), nil)
}

// NewJWTToken creates new JWT token.
func NewJWTToken(ja *jwtauth.JWTAuth, login string) (token jwt.Token, tokenString string, err error) {
    return ja.Encode(map[string]interface{}{string(UserLoginClaimName): login})
}

// NewCookieWithDefaults creates new cookie with defaults and parameter value.
func NewCookieWithDefaults(value string) *http.Cookie {
    return &http.Cookie{
        Name:     DefaultCookieName,
        Value:    value,
        Path:     DefaultCookiePath,
        MaxAge:   DefaultCookieMaxAge,
        SameSite: http.SameSiteDefaultMode,
    }
}

// HashPassword generates and returns hash of a given password.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash compares given password and hash.
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
