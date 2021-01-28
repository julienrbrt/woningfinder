package networking

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTemporaryErr struct {
	error
	temporary bool
}

func (e mockTemporaryErr) Temporary() bool {
	return e.temporary
}

func TestIsTemporaryError(t *testing.T) {
	a := assert.New(t)
	a.False(IsTemporaryError(errors.New("some random error")))
	a.False(IsTemporaryError(mockTemporaryErr{error: errors.New("some random error"), temporary: false}))
	a.True(IsTemporaryError(mockTemporaryErr{error: errors.New("some random error"), temporary: true}))
	a.False(IsTemporaryError(NewNetworkingError(&Response{StatusCode: http.StatusBadGateway}, "")))
	a.True(IsTemporaryError(NewTemporaryNetError(NewNetworkingError(&Response{StatusCode: http.StatusBadGateway}, ""))))
}

func TestIsNetworkingError(t *testing.T) {
	a := assert.New(t)
	a.False(IsNetworkingError(errors.New("some random error")))
	a.False(IsNetworkingError(mockTemporaryErr{error: errors.New("some random error"), temporary: false}))
	a.False(IsNetworkingError(mockTemporaryErr{error: errors.New("some random error"), temporary: true}))
	a.True(IsNetworkingError(NewNetworkingError(&Response{StatusCode: http.StatusBadGateway}, "")))
	a.True(IsNetworkingError(NewTemporaryNetError(NewNetworkingError(&Response{StatusCode: http.StatusBadGateway}, ""))))
}

func TestAsNetworkingError(t *testing.T) {
	a := assert.New(t)
	_, ok := AsNetworkingError(errors.New("some random error"))
	a.False(ok)
	_, ok = AsNetworkingError(mockTemporaryErr{error: errors.New("some random error"), temporary: false})
	a.False(ok)
	_, ok = AsNetworkingError(mockTemporaryErr{error: errors.New("some random error"), temporary: true})
	a.False(ok)
	mockErr := NewNetworkingError(&Response{StatusCode: http.StatusBadGateway}, "")
	netErr, ok := AsNetworkingError(mockErr)
	a.True(ok)
	a.Equal(mockErr, netErr)
	netErr, ok = AsNetworkingError(NewTemporaryNetError(mockErr))
	a.True(ok)
	a.Equal(mockErr, netErr)
}

func TestResponseFromError(t *testing.T) {
	a := assert.New(t)
	_, ok := ResponseFromError(errors.New("some random error"))
	a.False(ok)
	_, ok = ResponseFromError(mockTemporaryErr{error: errors.New("some random error"), temporary: false})
	a.False(ok)
	_, ok = ResponseFromError(mockTemporaryErr{error: errors.New("some random error"), temporary: true})
	a.False(ok)
	mockResponse := &Response{StatusCode: http.StatusBadGateway}
	resp, ok := ResponseFromError(NewNetworkingError(mockResponse, ""))
	a.True(ok)
	a.Equal(mockResponse, resp)
	resp, ok = ResponseFromError(NewTemporaryNetError(NewNetworkingError(mockResponse, "")))
	a.True(ok)
	a.Equal(mockResponse, resp)
}

func TestStatusCodeFromError(t *testing.T) {
	a := assert.New(t)
	_, ok := StatusCodeFromError(errors.New("some random error"))
	a.False(ok)
	_, ok = StatusCodeFromError(mockTemporaryErr{error: errors.New("some random error"), temporary: false})
	a.False(ok)
	_, ok = StatusCodeFromError(mockTemporaryErr{error: errors.New("some random error"), temporary: true})
	a.False(ok)
	statusCode, ok := StatusCodeFromError(NewNetworkingError(&Response{StatusCode: http.StatusBadGateway}, ""))
	a.True(ok)
	a.Equal(http.StatusBadGateway, statusCode)
	statusCode, ok = StatusCodeFromError(NewTemporaryNetError(NewNetworkingError(&Response{StatusCode: http.StatusInsufficientStorage}, "")))
	a.True(ok)
	a.Equal(http.StatusInsufficientStorage, statusCode)
}
