package networking

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/julienrbrt/woningfinder/pkg/networking/query"
)

// Request defines a request
type Request struct {
	Host      *url.URL
	Path      string
	Query     query.Query
	Method    string
	Headers   map[string]string
	Body      io.Reader
	bodyBytes []byte
	Context   context.Context
}

// Copy deep copies a request.
func (r *Request) Copy() (*Request, error) {
	copiedBody, err := r.CopyBody()
	if err != nil {
		return nil, fmt.Errorf("failed to copy request body, for creating a copy of the request")
	}

	var copiedBodyReader io.Reader
	if copiedBody != nil {
		copiedBodyReader = bytes.NewReader(copiedBody)
	}

	var copiedHeaders map[string]string
	if r.Headers != nil {
		copiedHeaders = make(map[string]string, len(r.Headers))
		for k, v := range r.Headers {
			copiedHeaders[k] = v
		}
	}

	var copiedHost *url.URL
	if r.Host != nil {
		copiedHostValue := *r.Host
		copiedHost = &copiedHostValue
	}

	var copiedQuery query.Query
	if r.Query != nil {
		copiedQuery = append(query.Query{}, r.Query...)
	}

	newRequest := *r
	newRequest.Host = copiedHost
	newRequest.Query = copiedQuery
	newRequest.Headers = copiedHeaders
	newRequest.Body = copiedBodyReader
	newRequest.bodyBytes = copiedBody

	return &newRequest, nil
}

// GetContext gets the request's context.
func (r *Request) GetContext() context.Context {
	if r.Context == nil {
		r.Context = context.Background()
	}

	return r.Context
}

// CopyBody reads the request body to bytes.
func (r *Request) CopyBody() ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	if r.bodyBytes == nil {
		if buf, ok := r.Body.(*bytes.Buffer); ok {
			r.bodyBytes = buf.Bytes()
		} else {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read request body for copying request %q: %w", r, err)
			}
			r.bodyBytes = b
		}
		r.Body = bytes.NewReader(r.bodyBytes)
	}

	return r.bodyBytes, nil
}

func (r *Request) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s %s", r.GetMethod(), r.Path)
}

// GetMethod gets the HTTP Method used for the request.
func (r *Request) GetMethod() string {
	if r.Method == "" {
		return http.MethodGet
	}
	return r.Method
}

// SetHeader sets an header in the request.
// key is the header name and value the value of the header.
func (r *Request) SetHeader(key, value string) {
	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
	r.Headers[key] = value
}
