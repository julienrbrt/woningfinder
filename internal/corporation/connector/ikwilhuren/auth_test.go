package ikwilhuren

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoginRequest(t *testing.T) {
	a := assert.New(t)
	req := loginRequest("username", "password")
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal("isAjax=true&password=password&username=username", string(body))
}
