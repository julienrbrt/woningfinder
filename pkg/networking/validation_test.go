package networking

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/julienrbrt/woningfinder/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestResponseValidator(t *testing.T) {
	expectations := map[int]bool{
		200: true,
		201: true,
		204: true,
		400: false,
		401: false,
		403: false,
		404: false,
		500: false,
		502: false,
		503: false,
	}

	a := assert.New(t)
	for statusCode, valid := range expectations {
		t.Run(fmt.Sprintf("status code: %d", statusCode), func(t *testing.T) {
			err := responseValidator(&Response{
				StatusCode: statusCode,
			})
			if valid {
				a.NoError(err)
			} else {
				a.Error(err)
			}
		})
	}
}

func TestDefaultResponseValidator_Validate_Temporary(t *testing.T) {
	expectations := map[int]bool{
		400: false,
		401: false,
		403: false,
		404: false,
		500: true,
		502: true,
		503: true,
	}

	a := assert.New(t)
	for statusCode, temporary := range expectations {
		t.Run(fmt.Sprintf("status code: %d", statusCode), func(t *testing.T) {
			err := responseValidator(&Response{
				StatusCode: statusCode,
			})
			var netErr net.Error
			a.True(errors.As(err, &netErr))
			a.Equal(temporary, netErr.Temporary())
		})
	}
}

func TestResponseValidator_FailingBodyReader(t *testing.T) {
	a := assert.New(t)
	bodyReaderErr := fmt.Errorf("error from the body reader")
	err := responseValidator(&Response{
		StatusCode: 502,
		Body:       util.NewNopCloser(util.NewErrReader(bodyReaderErr)),
	})
	a.Error(err)
	a.Contains(err.Error(), "502")
	a.Contains(err.Error(), bodyReaderErr.Error())
}
