package recovererr

import (
	"time"

	"github.com/cenkalti/backoff"
)

// ConstantBackoff implements backoff strategy using using exponentially increased delays.
type ExponentialBackoff struct {
	expBackoff backoff.ExponentialBackOff
	afterFunc  func(time.Duration) <-chan time.Time
}

// NewExponentialBackoff creates new exponential backoff using provided options.
func NewExponentialBackoff(opts ...Option) *ExponentialBackoff {
	bo := ExponentialBackoff{
		expBackoff: backoff.ExponentialBackOff{
			InitialInterval:     backoff.DefaultInitialInterval,
			RandomizationFactor: backoff.DefaultRandomizationFactor,
			Multiplier:          backoff.DefaultMultiplier,
			MaxInterval:         backoff.DefaultMaxInterval,
			MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
			Clock:               backoff.SystemClock,
		},
		afterFunc: time.After,
	}

	for _, opt := range opts {
		opt(&bo)
	}

	bo.expBackoff.Reset()

	return &bo
}

// Option configures exponential backoff parameters.
type Option func(*ExponentialBackoff)

// WithInitialInterval sets initial interval to exponential backoff.
func WithInitialInterval(initialInterval time.Duration) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.InitialInterval = initialInterval
	}
}

// WithMultiplier sets multiplier to exponential backoff.
func WithMultiplier(multiplier float64) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.Multiplier = multiplier
	}
}

// WithMaxInterval sets max interval to exponential backoff.
func WithMaxInterval(maxInterval time.Duration) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.MaxInterval = maxInterval
	}
}

// WithMaxElapsedTime sets max elapsed time to exponential backoff.
func WithMaxElapsedTime(maxElapsedTime time.Duration) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.MaxElapsedTime = maxElapsedTime
	}
}

// WithRandomisationFactory sets randomisation factor to exponential backoff.
func WithRandomisationFactory(randomisationFactor float64) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.RandomizationFactor = randomisationFactor
	}
}

// WithClock sets clock implementation to exponential backoff.
func WithClock(clock backoff.Clock) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.Clock = clock
	}
}

// Next implements the BackoffStrategy.Next method.
func (eb *ExponentialBackoff) Next() (time.Duration, bool) {
	d := eb.expBackoff.NextBackOff()
	if d == backoff.Stop {
		return 0, false
	}

	return d, true
}
