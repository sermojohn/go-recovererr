package recovererr

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// ConstantBackoff implements backoff strategy using using exponentially increased delays.
type ExponentialBackoff struct {
	impl backoff.ExponentialBackOff
}

// NewExponentialBackoff creates new exponential backoff using provided options.
func NewExponentialBackoff(opts ...ExponentialBackoffOption) *ExponentialBackoff {
	bo := ExponentialBackoff{
		impl: *backoff.NewExponentialBackOff(),
	}

	// customise defaults
	bo.impl.InitialInterval = 500 * time.Millisecond
	bo.impl.RandomizationFactor = 0
	bo.impl.Multiplier = 1.5
	bo.impl.MaxInterval = 10 * time.Second
	bo.impl.MaxElapsedTime = 30 * time.Second
	bo.impl.Clock = backoff.SystemClock

	for _, opt := range opts {
		opt(&bo)
	}

	bo.impl.Reset()

	return &bo
}

// ExponentialBackoffOption configures exponential backoff parameters.
type ExponentialBackoffOption func(*ExponentialBackoff)

// WithInitialInterval sets initial interval to exponential backoff.
func WithInitialInterval(initialInterval time.Duration) ExponentialBackoffOption {
	return func(eb *ExponentialBackoff) {
		eb.impl.InitialInterval = initialInterval
	}
}

// WithMultiplier sets multiplier to exponential backoff.
func WithMultiplier(multiplier float64) ExponentialBackoffOption {
	return func(eb *ExponentialBackoff) {
		eb.impl.Multiplier = multiplier
	}
}

// WithMaxInterval sets max interval to exponential backoff.
func WithMaxInterval(maxInterval time.Duration) ExponentialBackoffOption {
	return func(eb *ExponentialBackoff) {
		eb.impl.MaxInterval = maxInterval
	}
}

// WithMaxElapsedTime sets max elapsed time to exponential backoff.
func WithMaxElapsedTime(maxElapsedTime time.Duration) ExponentialBackoffOption {
	return func(eb *ExponentialBackoff) {
		eb.impl.MaxElapsedTime = maxElapsedTime
	}
}

// WithRandomisationFactory sets randomisation factor to exponential backoff.
func WithRandomisationFactory(randomisationFactor float64) ExponentialBackoffOption {
	return func(eb *ExponentialBackoff) {
		eb.impl.RandomizationFactor = randomisationFactor
	}
}

// WithClock sets clock implementation to exponential backoff.
func WithClock(clock backoff.Clock) ExponentialBackoffOption {
	return func(eb *ExponentialBackoff) {
		eb.impl.Clock = clock
	}
}

// Next implements the BackoffStrategy.Next method.
func (eb *ExponentialBackoff) Next() (time.Duration, bool) {
	d := eb.impl.NextBackOff()
	if d == backoff.Stop {
		return 0, false
	}

	return d, true
}
