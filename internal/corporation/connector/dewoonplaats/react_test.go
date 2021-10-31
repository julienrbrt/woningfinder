package dewoonplaats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReactRequest(t *testing.T) {
	a := assert.New(t)
	expected := `{"id":1,"method":"ReageerOpWoning","params":["testID"]}`
	req, err := reactRequest("testID")
	a.NoError(err)
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(expected, string(body))
}
