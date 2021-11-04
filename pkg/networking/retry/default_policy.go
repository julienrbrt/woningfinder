package retry

import (
	"time"
)

const (
	DefaultTimeout          = 12 * time.Second
	DefaultRetryInitialWait = 4 * time.Second
	DefaultRetryFactor      = 2
	DefaultRetryJitter      = time.Second
	DefaultRetryCount       = 3
)

// DefaultRetryPolicy returns a retry policy to be used in clients
func DefaultRetryPolicy() Policy {
	return ExponentialBackOff(DefaultRetryInitialWait, DefaultRetryFactor, DefaultRetryJitter, DefaultRetryCount)
}
