package networking

import (
	"errors"
)

type temporary interface {
	Temporary() bool
}

// IsTemporaryError checks if the given error is temporary.
func IsTemporaryError(err error) bool {
	var netErr temporary
	if errors.As(err, &netErr) {
		return netErr.Temporary()
	}
	return false
}

// IsNetworkingError checks if the given error is a networking error.
func IsNetworkingError(err error) bool {
	_, ok := AsNetworkingError(err)
	return ok
}

// AsNetworkingError type checks an error to a NetworkingError and if succeeded, it returns it
func AsNetworkingError(err error) (Error, bool) {
	var netErr Error
	if errors.As(err, &netErr) {
		return netErr, true
	}
	return nil, false
}

// ResponseFromError returns the error response.
func ResponseFromError(err error) (*Response, bool) {
	if netErr, ok := AsNetworkingError(err); ok {
		return netErr.Response(), true
	}
	return nil, false
}

// StatusCodeFromError returns the HTTP Status Code of the error.
func StatusCodeFromError(err error) (int, bool) {
	if response, ok := ResponseFromError(err); ok {
		return response.StatusCode, true
	}
	return 0, false
}
