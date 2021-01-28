package networking

import (
	"net"
)

// Error represents the networking error.
// Error inherits of net.Error and returns the last response.
type Error interface {
	net.Error
	Response() *Response
}

type networkingError struct {
	resp    *Response
	message string
}

// NewNetworkingError creates a networking error with a response and an error message.
func NewNetworkingError(resp *Response, message string) Error {
	return &networkingError{
		resp:    resp,
		message: message,
	}
}

func (e *networkingError) Timeout() bool {
	return false
}

func (e *networkingError) Temporary() bool {
	return false
}

func (e *networkingError) Response() *Response {
	return e.resp
}

func (e *networkingError) Error() string {
	return e.message
}

type temporaryNetError struct {
	wrapped net.Error
}

// NewTemporaryNetError creates an Error satisfying only net.Error (and not networking.Error)
func NewTemporaryNetError(err net.Error) net.Error {
	return temporaryNetError{wrapped: err}
}

func (e temporaryNetError) Error() string {
	return e.wrapped.Error()
}

func (e temporaryNetError) Timeout() bool {
	return e.wrapped.Timeout()
}

func (e temporaryNetError) Temporary() bool {
	return true
}

func (e temporaryNetError) Unwrap() error {
	return e.wrapped
}

type temporaryError struct {
	wrapped error
}

// NewTemporaryError returns a new temporary error
//
// Can be used in custom response validators, to designate a specific error
// as temporary for retry purposes.
func NewTemporaryError(err error) error {
	return &temporaryError{wrapped: err}
}

func (e *temporaryError) Error() string {
	return e.wrapped.Error()
}

func (e *temporaryError) Unwrap() error {
	return e.wrapped
}

func (e *temporaryError) Temporary() bool {
	return true
}
