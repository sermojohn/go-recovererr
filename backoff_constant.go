package recovererr

import "time"

// ConstantBackoff implements backoff strategy using constant delay and max attempts.
type ConstantBackoff struct {
	interval    time.Duration
	maxAttempts int

	attempt int
}

// NewConstantBackoff creates new constant backoff using provided parameters.
func NewConstantBackoff(opts ...ConstantBackoffOption) *ConstantBackoff {
	cb := ConstantBackoff{}

	for _, opt := range opts {
		opt(&cb)
	}

	if cb.interval == 0 {
		cb.interval = time.Second
	}
	if cb.maxAttempts == 0 {
		cb.maxAttempts = 3
	}

	return &cb
}

// ConstantBackoffOption configures constant backoff parameters.
type ConstantBackoffOption func(*ConstantBackoff)

// WithInterval configure constant backoff with specified backoff interval.
func WithInterval(d time.Duration) ConstantBackoffOption {
	return func(cb *ConstantBackoff) {
		cb.interval = d
	}
}

// WithMaxAttempts configure constant backoff with specified max backoff attempts.
func WithMaxAttempts(n int) ConstantBackoffOption {
	return func(cb *ConstantBackoff) {
		cb.maxAttempts = n
	}
}

// Next implements the BackoffStrategy.Next method.
func (cb *ConstantBackoff) Next() (time.Duration, bool) {
	cb.attempt++
	if cb.maxAttempts > 0 && cb.attempt > cb.maxAttempts {
		return 0, false
	}
	return cb.interval, true
}
