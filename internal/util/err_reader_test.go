package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrReader(t *testing.T) {
	a := assert.New(t)
	readerErr := errors.New("error from the error reader")
	r := NewErrReader(readerErr)
	a.Equal(readerErr, r.Error())
	_, err := r.Read(make([]byte, 0))
	a.Equal(err, readerErr)
}
