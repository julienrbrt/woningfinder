package networking_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/query"
	"github.com/woningfinder/woningfinder/pkg/util"
)

func TestRequestBuilder_NoHost(t *testing.T) {
	a := assert.New(t)
	_, err := networking.RequestBuilder(&networking.Request{})
	a.Error(err)
}

func TestRequestBuilder_InvalidBody(t *testing.T) {
	a := assert.New(t)
	readerErr := errors.New("error from the body reader")
	_, err := networking.RequestBuilder(&networking.Request{
		Host: &url.URL{Host: "example.com"},
		Body: ioutil.NopCloser(util.NewErrReader(readerErr)),
	})
	a.Error(err)
	a.True(errors.Is(err, readerErr))
}

func TestRequestBuilder_InvalidMethod(t *testing.T) {
	a := assert.New(t)
	_, err := networking.RequestBuilder(&networking.Request{
		Host:   &url.URL{Host: "example.com"},
		Method: "THIS IS INVALID",
	})
	a.Error(err)
}

func TestRequestBuilder_WithSlashes(t *testing.T) {
	a := assert.New(t)
	req, err := networking.RequestBuilder(&networking.Request{
		Host: &url.URL{Host: "example.com", Path: "a/path"},
		Path: "AH%2FWC P COACH A STD %2FQ %2FC",
	})
	a.NoError(err)
	a.Equal(http.MethodGet, req.Method)
	a.Equal("//example.com/a/path/AH%2FWC%20P%20COACH%20A%20STD%20%2FQ%20%2FC", req.URL.String())
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	a.NoError(err)
	a.Empty(body)
}

func TestRequestBuilder_SuccessMinimal(t *testing.T) {
	a := assert.New(t)
	req, err := networking.RequestBuilder(&networking.Request{
		Host: &url.URL{Host: "example.com"},
	})
	a.NoError(err)
	a.Equal(http.MethodGet, req.Method)
	a.Equal("//example.com", req.URL.String())
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	a.NoError(err)
	a.Empty(body)
}

func TestRequestBuilder_SuccessFull(t *testing.T) {
	a := assert.New(t)
	q := make(query.Query, 0, 2)
	q.Add("foo", "bar")
	q.Add("baz", "foobar")
	req, err := networking.RequestBuilder(&networking.Request{
		Host:   &url.URL{Scheme: "https", Host: "example.com", Path: "/base", RawQuery: "host=example"},
		Path:   "/api/v2/ping/",
		Method: http.MethodPost,
		Query:  q,
		Headers: map[string]string{
			"X-Custom-Header":         "ABC",
			"X-Custom-Header-Reverse": "CBA",
		},
		Body:    ioutil.NopCloser(strings.NewReader("request body")),
		Context: context.WithValue(context.Background(), "hello", "world"),
	})
	a.NoError(err)
	a.Equal(http.MethodPost, req.Method)
	a.Equal("https://example.com/base/api/v2/ping?foo=bar&baz=foobar", req.URL.String())
	a.Equal("ABC", req.Header.Get("X-Custom-Header"))
	a.Equal("CBA", req.Header.Get("X-Custom-Header-Reverse"))
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	a.NoError(err)
	a.Equal("request body", string(body))
	a.Equal("world", req.Context().Value("hello"))
}
