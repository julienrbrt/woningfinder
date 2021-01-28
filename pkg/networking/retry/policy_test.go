package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNonePolicy_BackOff(t *testing.T) {
	_, ok := None.BackOff(1)
	assert.False(t, ok)

	_, ok = None.BackOff(3)
	assert.False(t, ok)
}

func TestExponentialBackOffPolicy_BackOff(t *testing.T) {
	a := assert.New(t)

	p := ExponentialBackOff(500*time.Millisecond, 2, 200*time.Millisecond, 3)
	backOff, ok := p.BackOff(1)
	a.True(ok)
	a.GreaterOrEqual(int64(backOff), int64(300*time.Millisecond))
	a.LessOrEqual(int64(backOff), int64(700*time.Millisecond))

	backOff, ok = p.BackOff(3)
	a.True(ok)
	a.GreaterOrEqual(int64(backOff), int64(1800*time.Millisecond))
	a.LessOrEqual(int64(backOff), int64(2200*time.Millisecond))

	_, ok = p.BackOff(4)
	a.False(ok)
}

func TestExponentialBackOffPolicy_BackOffNegative(t *testing.T) {
	a := assert.New(t)

	p := ExponentialBackOff(-100*time.Millisecond, 1, 0, 1)
	backOff, ok := p.BackOff(1)
	a.True(ok)
	a.Equal(int64(0), int64(backOff))
}
