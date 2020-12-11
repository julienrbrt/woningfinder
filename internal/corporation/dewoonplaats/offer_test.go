package dewoonplaats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OfferRequest(t *testing.T) {
	a := assert.New(t)
	expected := `{"id":1,"method":"ZoekWoningen","params":[{"tehuur":true},"",true]}`
	req, err := offerRequest()
	a.NoError(err)
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(expected, string(body))
}
