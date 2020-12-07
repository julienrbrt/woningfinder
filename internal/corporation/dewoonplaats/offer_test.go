package dewoonplaats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OfferRequest(t *testing.T) {
	a := assert.New(t)
	expected := `{"id":1,"method":"ZoekWoningen","params":[{"prijsvanaf":500,"tehuur":true},"",true]}`
	req, err := offerRequest(500)
	a.NoError(err)
	body, err := req.CopyBody()
	a.NoError(err)
	a.Equal(expected, string(body))
}
