package dewoonplaats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoginRequest(t *testing.T) {
	a := assert.New(t)
	expected := `{"id":1,"method":"Login","params":["https://www.dewoonplaats.nl/mijn-woonplaats/","username","password",false]}`
	req, err := loginRequest("username", "password")
	a.NoError(err)
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(expected, string(body))
}
