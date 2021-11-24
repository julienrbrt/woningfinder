package middleware_test

import (
	"testing"

	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/julienrbrt/woningfinder/pkg/networking/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCreateDefaultHeadersMiddleware_AddHeaders(t *testing.T) {
	headers := map[string]string{
		"Content-Type":    "application-json",
		"Accept-Language": "en-US",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Charset":  "utf-8",
	}

	c := mockClient{}
	w := middleware.CreateDefaultHeadersMiddleware(headers)(&c)
	a := assert.New(t)

	req := &networking.Request{}
	resp, err := w.Send(req)

	a.EqualValues(headers, req.Headers)
	a.NoError(err)
	a.Len(resp.Request.Headers, 4)
}

func TestCreateDefaultHeadersMiddleware_OverrideDefaultHeaders(t *testing.T) {
	headers := map[string]string{
		"Content-Type":    "application-json",
		"Accept-Language": "en-US",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Charset":  "utf-8",
	}

	c := mockClient{}
	w := middleware.CreateDefaultHeadersMiddleware(headers)(&c)
	a := assert.New(t)

	req := &networking.Request{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp, err := w.Send(req)

	a.NoError(err)
	a.Len(resp.Request.Headers, 4)
	a.Equal(resp.Request.Headers["Content-Type"], "application/json")
}
