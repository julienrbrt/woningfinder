package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/woningfinder/woningfinder/pkg/networking"
	"github.com/woningfinder/woningfinder/pkg/networking/retry"
)

const MaxBackOffCount = 100

type policyContextKey struct{}

// ContextWithRetryPolicy adds a retryPolicy and its value to the given context
func ContextWithRetryPolicy(ctx context.Context, retryPolicy retry.Policy) context.Context {
	return context.WithValue(ctx, policyContextKey{}, retryPolicy)
}

// RetryCounterFromContext gets the number of request retries from a context
func RetryPolicyFromContext(ctx context.Context) (retry.Policy, bool) {
	v, ok := ctx.Value(policyContextKey{}).(retry.Policy)
	return v, ok
}

type counterContextKey struct{}

// ContextWithRetryCounter adds a retry counter value to the given context
func ContextWithRetryCounter(ctx context.Context, counter int) context.Context {
	return context.WithValue(ctx, counterContextKey{}, counter)
}

// RetryCounterFromContext gets the number of request retries from a context
func RetryCounterFromContext(ctx context.Context) int {
	if v, ok := ctx.Value(counterContextKey{}).(int); ok {
		return v
	}
	return 0
}

// CreateRetryMiddleware adds a retry policy to the request.
func CreateRetryMiddleware(defaultPolicy retry.Policy, sleeper retry.Sleeper) networking.ClientMiddleware {
	return func(base networking.Client) networking.Client {
		return networking.ClientFunc(func(origReq *networking.Request) (*networking.Response, error) {
			retryPolicy := getRetryPolicy(origReq.GetContext(), defaultPolicy)
			backOff, retryEnabled := retryPolicy.BackOff(1)

			var err error
			var resp *networking.Response
			for i := 1; i <= MaxBackOffCount; i++ {
				req, reqErr := getRequest(origReq, retryEnabled)
				if reqErr != nil {
					return nil, reqErr
				}

				resp, err = base.Send(req)
				if !retryEnabled || err == nil || !networking.IsTemporaryError(err) || isDeadlineBeforeBackOff(origReq.GetContext(), backOff) {
					break
				}

				sleeper(backOff)
				origReq.Context = incrementRetryCounter(origReq.GetContext())
				backOff, retryEnabled = retryPolicy.BackOff(i + 1)
			}

			return resp, err
		})
	}
}

func isDeadlineBeforeBackOff(ctx context.Context, backOff time.Duration) bool {
	deadline, ok := ctx.Deadline()
	return ok && deadline.Before(time.Now().Add(backOff))
}

func getRetryPolicy(ctx context.Context, def retry.Policy) retry.Policy {
	if p, ok := RetryPolicyFromContext(ctx); ok {
		return p
	}
	return def
}

func getRequest(req *networking.Request, retryEnabled bool) (*networking.Request, error) {
	if !retryEnabled {
		return req, nil
	}
	newReq, err := req.Copy()
	if err != nil {
		return nil, fmt.Errorf("failed to get a copy of the request of request %q for the retry policy: %w", req, err)
	}

	return newReq, nil
}

func incrementRetryCounter(ctx context.Context) context.Context {
	return ContextWithRetryCounter(ctx, RetryCounterFromContext(ctx)+1)
}
