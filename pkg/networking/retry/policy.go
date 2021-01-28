package retry

import (
	"math"
	"math/rand"
	"time"
)

// Policy defines a request retry policy.
// A Policy must implement a BackOff method, that will give the duration till the next retry and whether retries are allowed
type Policy interface {
	BackOff(n int) (time.Duration, bool)
}

// None means that there is no policy, and thus no retry of the request.
var None Policy = nonePolicy{}

type nonePolicy struct{}

func (nonePolicy) BackOff(_ int) (time.Duration, bool) {
	return 0, false
}

type exponentialBackOffPolicy struct {
	initialWait time.Duration
	factor      float64
	jitter      time.Duration
	count       int
}

// ExponentialBackOff is a policy that will retry the request by waiting a duration that will increase exponentially (if the request cannot be completed)
// The jitter argument permits to spread the requests and avoid to overload the server with too many retrial requests
func ExponentialBackOff(initialWait time.Duration, factor float64, jitter time.Duration, count int) Policy {
	return &exponentialBackOffPolicy{
		initialWait: initialWait,
		factor:      factor,
		jitter:      jitter,
		count:       count,
	}
}

func (p *exponentialBackOffPolicy) BackOff(n int) (time.Duration, bool) {
	if n > p.count {
		return 0, false
	}

	backOff := time.Duration(float64(p.initialWait) * math.Pow(p.factor, float64(n-1)))
	jitter := time.Duration(float64(p.jitter) * (rand.Float64() - 0.5) * 2)

	total := backOff + jitter
	if total < 0 {
		total = 0
	}

	return total, true
}
