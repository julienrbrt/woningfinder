package zig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoginRequest(t *testing.T) {
	a := assert.New(t)
	expected := `__hash__=hash&__id__=Account_Form_LoginFrontend&password=password&username=username`
	req := loginRequest("username", "password", "hash")
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(expected, string(body))
}
