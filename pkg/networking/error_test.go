package networking

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkingError_Error_Temporary(t *testing.T) {
	a := assert.New(t)
	m := "error-message"
	r := Response{StatusCode: 502}
	e := NewTemporaryNetError(NewNetworkingError(&r, m))
	a.Equal("error-message", e.Error())
	a.True(e.Temporary())
	a.False(e.Timeout())

	var networkingError Error
	a.True(errors.As(e, &networkingError))
	a.Equal(&r, networkingError.Response())
}

func TestNetworkingError_Error_NotTemporary(t *testing.T) {
	a := assert.New(t)
	m := "error-message"
	r := Response{StatusCode: 503}
	e := NewNetworkingError(&r, m)
	a.Equal("error-message", e.Error())
	a.False(e.Temporary())
	a.False(e.Timeout())
	a.Equal(&r, e.Response())
}

func TestNewTemporaryError_Error(t *testing.T) {
	a := assert.New(t)
	mockErr := fmt.Errorf("error-message")
	e := NewTemporaryError(mockErr)
	a.Equal("error-message", e.Error())
	a.True(IsTemporaryError(e))
	a.True(errors.Is(e, mockErr))
}
