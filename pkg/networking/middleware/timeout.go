package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/util"
)

const DefaultRequestTimeout = 25 * time.Second

// CreateTimeoutMiddleware sets a timeout to a request
func CreateTimeoutMiddleware(defaultTimeout time.Duration) networking.ClientMiddleware {
	return func(base networking.Client) networking.Client {
		return networking.ClientFunc(func(r *networking.Request) (*networking.Response, error) {
			ctx, cancelFunc := context.WithDeadline(r.GetContext(), time.Now().Add(defaultTimeout))
			r.Context = ctx

			resp, err := base.Send(r)
			if err != nil {
				return nil, addCancelFuncToErr(err, cancelFunc)
			}

			return addCancelFuncToResponse(resp, cancelFunc), nil
		})
	}
}

func addCancelFuncToErr(err error, cancelFunc context.CancelFunc) error {
	var netErr networking.Error
	if errors.As(err, &netErr) {
		_ = addCancelFuncToResponse(netErr.Response(), cancelFunc)
	}

	return err
}

func addCancelFuncToResponse(resp *networking.Response, cancelFunc context.CancelFunc) *networking.Response {
	resp.SetBody(util.NewCancelFuncReadCloser(resp.Body, cancelFunc))
	return resp
}
