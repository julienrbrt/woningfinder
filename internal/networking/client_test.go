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
	"github.com/woningfinder/woningfinder/internal/networking"
)

type mockHTTPClient struct {
	err         error
	statusCode  int
	response    *http.Response
	lastRequest *http.Request
}

func (c *mockHTTPClient) Do(r *http.Request) (*http.Response, error) {
	c.lastRequest = r
	if c.err != nil {
		return nil, c.err
	}

	statusCode := 200
	if c.statusCode > 0 {
		statusCode = c.statusCode
	}

	resp := c.response
	if resp == nil {
		resp = &http.Response{
			StatusCode: statusCode,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(strings.NewReader("mock-response")),
		}
	}

	return resp, nil
}

func mockMiddleware(c networking.Client) networking.Client {
	return networking.ClientFunc(func(r *networking.Request) (*networking.Response, error) {
		r.SetHeader("X-Mock-Middleware", "true")
		return c.Send(r)
	})
}

func TestClient_Send_FailingHttpClient(t *testing.T) {
	a := assert.New(t)
	httpClient := &mockHTTPClient{err: errors.New("mock-error")}
	c := networking.NewClient(httpClient)

	_, err := c.Send(&networking.Request{Host: &url.URL{Host: "www.example.com"}})
	a.True(errors.Is(errors.Unwrap(err), httpClient.err))
}

func TestClient_Send_Success(t *testing.T) {
	a := assert.New(t)
	httpClient := &mockHTTPClient{}
	c := networking.NewClient(httpClient, mockMiddleware)

	ctx := context.WithValue(context.Background(), "a", "b")
	req := &networking.Request{
		Host:    &url.URL{Scheme: "https", Host: "www.example.com"},
		Method:  http.MethodPost,
		Path:    "Test",
		Body:    strings.NewReader("request-content"),
		Headers: map[string]string{"X-Foo": "Bar"},
		Context: ctx,
	}
	resp, err := c.Send(req)
	a.NoError(err)
	a.Equal(req, resp.Request)
	a.NotNil(resp.RawResponse)
	body, _ := resp.CopyBody()
	a.Equal("mock-response", string(body))
}

func createMockClientFunc(err error) networking.ClientFunc {
	return func(r *networking.Request) (*networking.Response, error) {
		if err != nil {
			return nil, err
		}
		return &networking.Response{
			Request: r,
		}, nil
	}
}

func TestClientFunc_Send_Success(t *testing.T) {
	a := assert.New(t)
	r := &networking.Request{}
	resp, err := createMockClientFunc(nil).Send(r)
	a.NoError(err)
	a.Equal(r, resp.Request)
}

func TestClientFunc_Send_Failure(t *testing.T) {
	a := assert.New(t)
	clientErr := errors.New("client-error")
	_, err := createMockClientFunc(clientErr).Send(&networking.Request{})
	a.Error(err)
	a.True(errors.Is(err, clientErr))
}

func TestClient_Send_SuccessResponseValidator(t *testing.T) {
	a := assert.New(t)
	httpClient := &mockHTTPClient{response: &http.Response{
		StatusCode: 200,
	}}
	c := networking.NewClient(httpClient)

	_, err := c.Send(&networking.Request{
		Host: &url.URL{Host: "www.example.com"},
		Path: "Test",
		Body: strings.NewReader("request-content"),
	})
	a.Nil(err)
}

func TestClient_Send_NotFound_SuccessResponseValidator(t *testing.T) {
	a := assert.New(t)
	httpClient := &mockHTTPClient{response: &http.Response{
		StatusCode: 404,
	}}
	c := networking.NewClient(httpClient)

	_, err := c.Send(&networking.Request{
		Host: &url.URL{Host: "www.example.com"},
		Path: "Test",
		Body: strings.NewReader("request-content"),
	})
	a.Error(err)
}
