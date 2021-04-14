package errors

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse defines the response returned when an error happens
type ErrorResponse struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

var (
	ErrMethodNotAllowed = &ErrorResponse{StatusCode: http.StatusMethodNotAllowed, Message: "Method not allowed"}
	ErrNotFound         = &ErrorResponse{StatusCode: http.StatusNotFound, Message: "Resource not found"}
	ErrBadRequest       = &ErrorResponse{StatusCode: http.StatusBadRequest, Message: "Bad request"}
	ErrUnauthorized     = &ErrorResponse{StatusCode: http.StatusUnauthorized, StatusText: "Unauthorized"}
)

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func ErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}

func ServerErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		StatusText: "Internal server error",
		Message:    err.Error(),
	}
}
