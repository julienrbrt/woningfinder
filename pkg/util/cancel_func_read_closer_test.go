package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cancelFuncHandler struct {
	cancelled bool
}

func (h *cancelFuncHandler) cancel() {
	h.cancelled = true
}

func TestCancelFuncReadCloser(t *testing.T) {
	a := assert.New(t)
	c1 := NewNopCloser(strings.NewReader("test123"))
	h := cancelFuncHandler{}
	c2 := NewCancelFuncReadCloser(c1, h.cancel)

	count, err := c2.Read(make([]byte, 2))
	a.Equal(2, count)
	a.NoError(err)
	a.False(c1.Closed())
	a.False(h.cancelled)

	a.NoError(c2.Close())
	a.True(c1.Closed())
	a.True(h.cancelled)
}
