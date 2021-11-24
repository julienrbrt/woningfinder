package networking

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"testing"

	"github.com/julienrbrt/woningfinder/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestResponse_CopyBody(t *testing.T) {
	a := assert.New(t)
	bodyReader := util.NewNopCloser(strings.NewReader("response-content"))
	r := Response{
		Body: bodyReader,
	}
	copiedBody, err := r.CopyBody()
	a.NoError(err)
	a.Equal("response-content", string(copiedBody))
	a.True(bodyReader.Closed())

	originalBody, err := r.CopyBody()
	a.NoError(err)
	a.Equal("response-content", string(originalBody))
}

func TestResponse_CopyBodyNil(t *testing.T) {
	a := assert.New(t)
	r := Response{}
	copiedBody, err := r.CopyBody()
	a.NoError(err)
	a.Equal("", string(copiedBody))
}

func TestResponse_CopyBody_Failure(t *testing.T) {
	a := assert.New(t)
	mockErr := errors.New("mock-error")
	bodyReader := util.NewNopCloser(util.NewErrReader(mockErr))
	r := Response{
		Body: bodyReader,
	}
	_, err := r.CopyBody()
	a.Error(err)
	a.True(errors.Is(err, mockErr))
	a.True(bodyReader.Closed())
}

func TestResponse_ReadJSONBody_EmptyBody(t *testing.T) {
	a := assert.New(t)
	r := Response{}

	var data map[string]interface{}
	a.Error(r.ReadJSONBody(&data))
}

func TestResponse_ReadJSONBody_FailingReader(t *testing.T) {
	a := assert.New(t)
	bodyReader := util.NewNopCloser(util.NewErrReader(fmt.Errorf("error from the reader")))
	r := Response{
		Body: bodyReader,
	}

	var data map[string]interface{}
	a.Error(r.ReadJSONBody(&data))
	a.True(bodyReader.Closed())
}

func TestResponse_ReadJSONBody_InvalidJSON(t *testing.T) {
	a := assert.New(t)
	bodyReader := util.NewNopCloser(strings.NewReader("invalid json"))
	r := Response{
		Body: bodyReader,
	}

	var data map[string]interface{}
	a.Error(r.ReadJSONBody(&data))
	a.True(bodyReader.Closed())
}

func TestResponse_ReadJSONBody_Success(t *testing.T) {
	a := assert.New(t)
	bodyReader := util.NewNopCloser(strings.NewReader(`{"foo": "bar"}`))
	r := Response{
		Body: bodyReader,
	}

	var data map[string]interface{}
	a.NoError(r.ReadJSONBody(&data))
	foo, ok := data["foo"].(string)
	a.True(ok)
	a.Equal("bar", foo)
	a.True(bodyReader.Closed())
	a.NoError(r.ReadJSONBody(&data))
}

func TestResponse_SetBody_NoRawResponse(t *testing.T) {
	a := assert.New(t)
	r := Response{}
	r.SetBody(util.NewNopCloser(strings.NewReader("test123")))
	a.Nil(r.RawResponse)
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	a.NoError(err)
	a.Equal("test123", string(body))
}

func TestResponse_SetBody_WithRawResponse(t *testing.T) {
	a := assert.New(t)
	r := Response{RawResponse: &http.Response{}}
	r.SetBody(util.NewNopCloser(strings.NewReader("test123")))
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	a.NoError(err)
	a.Equal("test123", string(body))
}

func TestResponse_CopyBodyThenLogRawResponseBody(t *testing.T) {
	a := assert.New(t)
	body := ioutil.NopCloser(strings.NewReader("test123"))
	r := Response{
		Body:        body,
		RawResponse: &http.Response{Body: body},
	}
	copiedBody, err := r.CopyBody()
	a.NoError(err)
	dumpedBody, err := httputil.DumpResponse(r.RawResponse, true)
	a.NoError(err)
	a.Equal("test123", string(copiedBody))
	a.Contains(string(dumpedBody), "test123")
}
