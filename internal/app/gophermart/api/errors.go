package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

var (
	ErrBadRequest          = &ErrorResponse{StatusCode: 400, Message: "Bad request"}
	ErrWrongLoginPassword  = &ErrorResponse{StatusCode: 401, Message: "Wrong login/password"}
	ErrLoginIsAlreadyTaken = &ErrorResponse{StatusCode: 409, Message: "Login has already been taken"}
	ErrInvalidOrderNumber  = &ErrorResponse{StatusCode: 422, Message: "Invalid order number"}
)

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}
func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}
func ServerErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 500,
		StatusText: "Internal server error",
		Message:    err.Error(),
	}
}
