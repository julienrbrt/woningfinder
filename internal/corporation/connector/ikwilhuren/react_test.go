package ikwilhuren

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReactRequest(t *testing.T) {
	a := assert.New(t)
	req := reactRequest("testID")
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(`day=&object=testID&time=`, string(body))
}
