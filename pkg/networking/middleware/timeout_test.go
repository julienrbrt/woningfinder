package middleware_test

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/middleware"
	"github.com/woningfinder/woningfinder/pkg/util"
)

func TestCreateTimeoutMiddleware_Success(t *testing.T) {
	a := assert.New(t)
	c := mockClient{}
	w := middleware.CreateTimeoutMiddleware(42 * time.Second)(&c)

	req := &networking.Request{}
	resp, err := w.Send(req)
	a.NoError(err)
	deadline, ok := resp.Request.GetContext().Deadline()
	a.True(ok)
	a.True(deadline.After(time.Now()))
	a.True(deadline.Before(time.Now().Add(43 * time.Second)))
	a.NoError(resp.Request.Context.Err())

	_, err = util.ReadAllAndClose(resp.Body)
	a.NoError(err)
	a.Error(resp.Request.Context.Err())
}

func TestCreateTimeoutMiddleware_Skip(t *testing.T) {
	a := assert.New(t)
	c := mockClient{}
	w := middleware.CreateTimeoutMiddleware(42 * time.Second)(&c)

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(12*time.Second))
	req := &networking.Request{
		Context: ctx,
	}
	resp, err := w.Send(req)
	a.NoError(err)
	deadline, ok := resp.Request.GetContext().Deadline()
	a.True(ok)
	a.True(deadline.After(time.Now()))
	a.True(deadline.Before(time.Now().Add(13 * time.Second)))
}

func TestCreateTimeoutMiddleware_Failure(t *testing.T) {
	a := assert.New(t)
	c := mockClient{err: errors.New("mock-error")}
	w := middleware.CreateTimeoutMiddleware(19 * time.Second)(&c)

	req := &networking.Request{}
	_, err := w.Send(req)
	a.Error(err)
	a.True(errors.Is(err, c.err))
}

func TestCreateTimeoutMiddleware_NetworkingError(t *testing.T) {
	a := assert.New(t)
	req := &networking.Request{}
	c := mockClient{err: networking.NewNetworkingError(&networking.Response{
		Request: req,
		Body:    ioutil.NopCloser(strings.NewReader("test123"))}, ""),
	}
	w := middleware.CreateTimeoutMiddleware(19 * time.Second)(&c)
	_, err := w.Send(req)
	a.Error(err)
	a.True(errors.Is(err, c.err))

	var netErr networking.Error
	a.True(errors.As(err, &netErr))
	a.NoError(netErr.Response().Request.Context.Err())
	_, err = util.ReadAllAndClose(netErr.Response().Body)
	a.NoError(err)
	a.Error(netErr.Response().Request.Context.Err())
}
