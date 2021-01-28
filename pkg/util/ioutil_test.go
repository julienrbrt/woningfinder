package util

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadAllAndClose_FailingReader(t *testing.T) {
	a := assert.New(t)
	r := NewNopCloser(NewErrReader(errors.New("error from the reader")))
	_, err := ReadAllAndClose(r)
	a.Error(err)
	a.True(r.Closed())
}

func TestReadAllAndClose_Success(t *testing.T) {
	a := assert.New(t)
	r := NewNopCloser(strings.NewReader("test123"))
	contents, err := ReadAllAndClose(r)
	a.NoError(err)
	a.Equal("test123", string(contents))
	a.True(r.Closed())
}
