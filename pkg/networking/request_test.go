package networking

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/pkg/networking/query"
)

func TestRequest_GetMethodDefault(t *testing.T) {
	r := Request{}
	if r.GetMethod() != http.MethodGet {
		t.Errorf("Expected method to be GET, got %s", r.GetMethod())
	}
}

func TestRequest_GetMethod(t *testing.T) {
	r := Request{Method: http.MethodHead}
	if r.GetMethod() != http.MethodHead {
		t.Errorf("Expected method to be HEAD, got %s", r.GetMethod())
	}
}

func TestRequest_CopyBody(t *testing.T) {
	r := Request{
		Body: strings.NewReader("request-content"),
	}

	body, err := r.CopyBody()
	if err != nil {
		t.Fatalf("Expected CopyBody not to return an error, got %s", err)
	}
	if string(body) != "request-content" {
		t.Errorf("Expected the request body to be 'request-content', got %s", string(body))
	}

	originalBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Expected read all from original request body reader not to return an error, got %s", err)
	}
	if string(originalBody) != "request-content" {
		t.Errorf("Expected the original request body to be 'request-content', got %s", string(originalBody))
	}
}

func TestRequest_CopyBody_BytesBuffer(t *testing.T) {
	r := Request{
		Body: bytes.NewBuffer([]byte("request-content")),
	}

	body, err := r.CopyBody()
	if err != nil {
		t.Fatalf("Expected CopyBody not to return an error, got %s", err)
	}
	if string(body) != "request-content" {
		t.Errorf("Expected the request body to be 'request-content', got %s", string(body))
	}

	originalBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Expected read all from original request body reader not to return an error, got %s", err)
	}
	if string(originalBody) != "request-content" {
		t.Errorf("Expected the original request body to be 'request-content', got %s", string(originalBody))
	}
}

func TestRequest_CopyBodyNil(t *testing.T) {
	r := Request{}
	body, err := r.CopyBody()
	if err != nil {
		t.Fatalf("Expected no error form CopyBody when the body is nil, got %s", err)
	}
	if body != nil {
		t.Errorf("Expected the response body to be nil, got %q", string(body))
	}
}

func TestRequest_Copy_Minimal(t *testing.T) {
	a := assert.New(t)
	r := &Request{}
	r2, err := r.Copy()
	a.NoError(err)
	a.Equal(r, r2)
	a.False(r == r2)

	r2.Host = &url.URL{Scheme: "http"}
	r2.Query = query.Query{query.Pair{Key: "baz", Value: "foobar"}}
	r2.Headers = map[string]string{"X-Custom": "Foo"}
	r2.Context = context.WithValue(context.Background(), "bar", "baz")

	a.Nil(r.Host)
	a.Empty(r.Path)
	a.Nil(r.Query)
	a.Nil(r.Headers)
	a.Nil(r.Context)
}

func TestRequest_Copy_Full(t *testing.T) {
	a := assert.New(t)
	r := &Request{
		Host:    &url.URL{Host: "example.com"},
		Path:    "/alias/1",
		Query:   query.Query{query.Pair{Key: "foo", Value: "bar"}},
		Method:  "POST",
		Headers: map[string]string{"X-Custom": "Foobar"},
		Body:    ioutil.NopCloser(strings.NewReader("test123")),
		Context: context.WithValue(context.Background(), "foo", "bar"),
	}
	r2, err := r.Copy()
	a.NoError(err)
	a.Equal(r, r2)
	a.False(r == r2)

	rBody, err := ioutil.ReadAll(r.Body)
	a.NoError(err)
	r2Body, err := ioutil.ReadAll(r2.Body)
	a.NoError(err)
	a.Equal(rBody, r2Body)

	r2.Host.Scheme = "http"
	r2.Query.Add("baz", "foobar")
	r2.Headers["X-Custom2"] = "Bar"
	r2.Context = context.WithValue(r2.Context, "bar", "baz")

	a.Empty(r.Host.Scheme)
	a.Equal("/alias/1", r.Path)
	a.Len(r.Query, 1)
	a.Equal("foo", r.Query[0].Key)
	a.Equal("bar", r.Query[0].Value)
	a.Len(r.Headers, 1)
	a.Equal("Foobar", r.Headers["X-Custom"])
	val, ok := r.Context.Value("foo").(string)
	a.True(ok)
	a.Equal("bar", val)
	_, ok = r.Context.Value("bar").(string)
	a.False(ok)
}
