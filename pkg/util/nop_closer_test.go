package util

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNopCloser(t *testing.T) {
	a := assert.New(t)
	c := NewNopCloser(strings.NewReader("test123"))

	count, err := c.Read(make([]byte, 2))
	a.Equal(2, count)
	a.NoError(err)
	a.False(c.Closed())

	a.NoError(c.Close())
	a.True(c.Closed())

	count, err = c.Read(make([]byte, 2))
	a.Equal(0, count)
	a.Error(err)
	a.True(c.Closed())
}
