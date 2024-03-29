package util_test

import (
	"testing"

	"github.com/julienrbrt/woningfinder/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestEmailCheck(t *testing.T) {
	a := assert.New(t)
	a.False(util.IsEmailValid("foo bar.com"))
	a.False(util.IsEmailValid(""))
	a.True(util.IsEmailValid("f@example-domain.org"))
	a.True(util.IsEmailValid("foo+spam@bar.com"))
	a.True(util.IsEmailValid("foo@bar.com"))
}
