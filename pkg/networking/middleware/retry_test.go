package middleware_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/julienrbrt/woningfinder/pkg/networking"
	"github.com/julienrbrt/woningfinder/pkg/networking/middleware"
	"github.com/julienrbrt/woningfinder/pkg/networking/query"
	"github.com/julienrbrt/woningfinder/pkg/networking/retry"
	"github.com/stretchr/testify/assert"
)

type mockSleeper struct {
	sleeps []time.Duration
}

func (s *mockSleeper) sleep(t time.Duration) {
	s.sleeps = append(s.sleeps, t)
}

type infiniteRetryPolicy time.Duration

func (p infiniteRetryPolicy) BackOff(_ int) (time.Duration, bool) {
	return time.Duration(p), true
}

type mockRetryPolicy []time.Duration

func (p mockRetryPolicy) BackOff(n int) (time.Duration, bool) {
	if n > len(p) {
		return 0, false
	}

	return p[n-1], true
}

func TestCreateRetryMiddleware_NoPolicy(t *testing.T) {
	req := networking.Request{
		Method: "POST",
		Path:   "/test123",
		Body:   strings.NewReader("test-content"),
		Query: query.Query{
			query.Pair{Key: "a", Value: "b"},
		},
	}
	c := mockClient{
		err: networking.NewTemporaryNetError(networking.NewNetworkingError(&networking.Response{StatusCode: 502}, "")),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(retry.None, s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 1)
	a.Equal(0, middleware.RetryCounterFromContext(c.lastRequests[0].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Equal("test-content", string(r.body))
	}
}

func TestCreateRetryMiddleware_NoNetError(t *testing.T) {
	req := networking.Request{
		Context: middleware.ContextWithRetryPolicy(context.Background(), mockRetryPolicy{1 * time.Second, 2 * time.Second}),
		Method:  "POST",
		Path:    "/test123",
		Body:    strings.NewReader("test-content"),
		Query: query.Query{
			query.Pair{Key: "a", Value: "b"},
		},
	}
	c := mockClient{
		err: errors.New("this is not a net.Error"),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(retry.None, s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 1)
	for _, r := range c.lastRequests {
		a.Equal("test-content", string(r.body))
	}
}

func TestCreateRetryMiddleware_NoTemporaryError(t *testing.T) {
	req := networking.Request{
		Context: middleware.ContextWithRetryPolicy(context.Background(), mockRetryPolicy{1 * time.Second, 2 * time.Second}),
		Method:  "POST",
		Path:    "/test123",
		Body:    strings.NewReader("test-content"),
		Query: query.Query{
			query.Pair{Key: "a", Value: "b"},
		},
	}
	c := mockClient{
		err: networking.NewNetworkingError(&networking.Response{StatusCode: 404}, ""),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(retry.None, s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 1)
	a.Equal(0, middleware.RetryCounterFromContext(c.lastRequests[0].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Equal("test-content", string(r.body))
	}
}

func TestCreateRetryMiddleware_WithBody(t *testing.T) {
	backOffs := []time.Duration{1 * time.Second, 2 * time.Second}
	req := networking.Request{
		Context: middleware.ContextWithRetryPolicy(context.Background(), mockRetryPolicy(backOffs)),
		Method:  "POST",
		Path:    "/test123",
		Body:    strings.NewReader("test-content"),
		Query: query.Query{
			query.Pair{Key: "a", Value: "b"},
		},
	}
	c := mockClient{
		err: networking.NewTemporaryNetError(networking.NewNetworkingError(&networking.Response{StatusCode: 502}, "")),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(retry.None, s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 3)
	a.Equal(2, middleware.RetryCounterFromContext(c.lastRequests[2].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Equal("test-content", string(r.body))
	}
}

func TestCreateRetryMiddleware_NoBody(t *testing.T) {
	backOffs := []time.Duration{1 * time.Second, 2 * time.Second}
	req := networking.Request{
		Context: middleware.ContextWithRetryPolicy(context.Background(), mockRetryPolicy(backOffs)),
		Method:  "GET",
		Path:    "/test123",
		Query: query.Query{
			query.Pair{Key: "a", Value: "b"},
		},
	}
	c := mockClient{
		err: networking.NewTemporaryNetError(networking.NewNetworkingError(&networking.Response{StatusCode: 502}, "")),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(retry.None, s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 3)
	a.Equal(2, middleware.RetryCounterFromContext(c.lastRequests[2].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Nil(r.body)
	}
}

func TestCreateRetryMiddleware_Infinite(t *testing.T) {
	req := networking.Request{
		Context: middleware.ContextWithRetryPolicy(context.Background(), infiniteRetryPolicy(1*time.Second)),
		Method:  "POST",
		Path:    "/test123",
		Body:    strings.NewReader("test-content"),
	}
	c := mockClient{
		err: networking.NewTemporaryNetError(networking.NewNetworkingError(&networking.Response{StatusCode: 502}, "")),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(retry.None, s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Equal(middleware.MaxBackOffCount-1, middleware.RetryCounterFromContext(c.lastRequests[middleware.MaxBackOffCount-1].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Equal("test-content", string(r.body))
	}
}

func TestCreateRetryMiddleware_Deadline(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelFunc()

	req := networking.Request{
		Context: ctx,
		Method:  "POST",
		Path:    "/test123",
		Body:    strings.NewReader("test-content"),
	}
	c := mockClient{
		err: networking.NewTemporaryNetError(networking.NewNetworkingError(&networking.Response{StatusCode: 502}, "")),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(infiniteRetryPolicy(3*time.Second), s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Equal(2, middleware.RetryCounterFromContext(c.lastRequests[2].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Equal("test-content", string(r.body))
	}
}

func TestCreateRetryMiddleware_Default(t *testing.T) {
	backOffs := []time.Duration{1 * time.Second, 2 * time.Second}
	req := networking.Request{
		Method: "GET",
		Path:   "/test123",
	}
	c := mockClient{
		err: networking.NewTemporaryNetError(networking.NewNetworkingError(&networking.Response{StatusCode: 502}, "")),
	}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(mockRetryPolicy(backOffs), s.sleep)(&c)
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 3)
	a.Equal(2, middleware.RetryCounterFromContext(c.lastRequests[2].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Nil(r.body)
	}
}

func TestCreateRetryMiddleware_ContextShouldNotBeReused(t *testing.T) {
	counterMiddleware := func(base networking.Client) networking.Client {
		return networking.ClientFunc(func(req *networking.Request) (*networking.Response, error) {
			counter, ok := req.GetContext().Value("counter").(int)
			if !ok {
				counter = 0
			}
			req.Context = context.WithValue(req.GetContext(), "counter", counter+1)

			return base.Send(req)
		})
	}

	backOffs := []time.Duration{1 * time.Second, 2 * time.Second}
	req := networking.Request{
		Method: "GET",
		Path:   "/test123",
	}
	c := mockClient{err: context.DeadlineExceeded}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(mockRetryPolicy(backOffs), s.sleep)(counterMiddleware(&c))
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 3)
	a.Equal(2, middleware.RetryCounterFromContext(c.lastRequests[2].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Equal(1, r.req.GetContext().Value("counter").(int))
	}
}

func TestCreateRetryMiddleware_RequestTimeout(t *testing.T) {
	backOffs := []time.Duration{1 * time.Second, 2 * time.Second}
	req := networking.Request{
		Method: "GET",
		Path:   "/test123",
	}
	c := mockClient{err: context.DeadlineExceeded}
	s := mockSleeper{}
	w := middleware.CreateRetryMiddleware(mockRetryPolicy(backOffs), s.sleep)(middleware.CreateTimeoutMiddleware(-1 * time.Second)(&c))
	_, err := w.Send(&req)

	a := assert.New(t)
	a.Error(err)
	a.True(errors.Is(err, c.err))
	a.Len(c.lastRequests, 3)
	a.Equal(2, middleware.RetryCounterFromContext(c.lastRequests[2].req.GetContext()))
	for _, r := range c.lastRequests {
		a.Nil(r.body)
	}
}
