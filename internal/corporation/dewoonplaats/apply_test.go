package dewoonplaats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ApplyRequest(t *testing.T) {
	a := assert.New(t)
	expected := `{"id":1,"method":"ReageerOpWoning","params":["testID"]}`
	req, err := applyRequest("testID")
	a.NoError(err)
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(expected, string(body))
}