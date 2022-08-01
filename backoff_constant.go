package recovererr

import "time"

// ConstantBackoff implements backoff strategy using constant delay and max attempts.
type ConstantBackoff struct {
	delay     time.Duration
	afterFunc func(time.Duration) <-chan time.Time
	max       int

	counter int
}

// NewConstantBackoff creates new constant backoff using provided parameters.
func NewConstantBackoff(d time.Duration, max int) *ConstantBackoff {
	cb := ConstantBackoff{
		delay:     d,
		afterFunc: time.After,
		max:       max,
	}
	return &cb
}

// Next implements the BackoffStrategy.Next method.
func (cb *ConstantBackoff) Next() (time.Duration, bool) {
	cb.counter++
	if cb.max > 0 && cb.counter > cb.max {
		return 0, false
	}
	return cb.delay, true
}
