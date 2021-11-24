package middleware_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/julienrbrt/woningfinder/pkg/networking/middleware"

	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/stretchr/testify/assert"
)

func TestCreateHostMiddleware_Success(t *testing.T) {
	a := assert.New(t)
	c := mockClient{}
	h := &url.URL{Scheme: "https", Host: "example.com"}
	w := middleware.CreateHostMiddleware(&url.URL{Scheme: "https", Host: "example.com"})(&c)

	_, err := w.Send(&networking.Request{})
	a.NoError(err)
	a.Len(c.lastRequests, 1)
	a.Equal(h, c.lastRequests[0].req.Host)
}

func TestCreateHostMiddleware_Failure(t *testing.T) {
	a := assert.New(t)
	c := mockClient{err: errors.New("error from the mock client")}
	h := &url.URL{Scheme: "https", Host: "example.com"}
	w := middleware.CreateHostMiddleware(&url.URL{Scheme: "https", Host: "example.com"})(&c)

	_, err := w.Send(&networking.Request{})
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 1)
	a.Equal(h, c.lastRequests[0].req.Host)
}
