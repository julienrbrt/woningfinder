package networking

import (
	"fmt"
	"net/http"
)

// BaseClient represents a client that can perform HTTP calls.
// A Client can be built on top of any BaseClient (see NewClient).
type BaseClient interface {
	Do(r *http.Request) (*http.Response, error)
}

// Client defines a client that sends a Request and returns a Response
type Client interface {
	Send(r *Request) (*Response, error)
}

// ClientFunc type to permit to have the Client's Send signature its own type.
// ClientFunc implements the Client interface.
type ClientFunc func(r *Request) (*Response, error)

// Send can be used on ClientFunc in order to send the request and get the response.
func (f ClientFunc) Send(r *Request) (*Response, error) {
	return f(r)
}

// ClientMiddleware is a middleware permitting to pass clients.
type ClientMiddleware func(Client) Client

type client struct {
	baseClient BaseClient
}

// NewClient permits to create a networking client
// A client can use ClientMiddleware in order to enrich the client with multiple middlewares.
// The middlewares must be sorted from the first to be passed through to the last: e.g. for middleware [m1, m2] we have client -> m1 -> m2 -> ... -> client.
func NewClient(baseClient BaseClient, middleware ...ClientMiddleware) Client {
	var c Client = &client{
		baseClient: baseClient,
	}
	for i := len(middleware) - 1; i >= 0; i-- {
		c = middleware[i](c)
	}
	return c
}

func (c *client) Send(r *Request) (*Response, error) {
	req, err := RequestBuilder(r)
	if err != nil {
		return nil, err
	}

	resp, err := c.baseClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed sending request: %s %s: %w", req.Method, req.URL, err)
	}

	response := &Response{
		Request:     r,
		StatusCode:  resp.StatusCode,
		Headers:     resp.Header,
		Body:        resp.Body,
		RawRequest:  req,
		RawResponse: resp,
	}

	if err := responseValidator(response); err != nil {
		return nil, err
	}

	return response, nil
}
