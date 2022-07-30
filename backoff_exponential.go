package recovererr

import (
	"time"

	"github.com/cenkalti/backoff"
)

type ExponentialBackoff struct {
	expBackoff backoff.ExponentialBackOff
	afterFunc  func(time.Duration) <-chan time.Time
}

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

type Option func(*ExponentialBackoff)

func WithInitialInterval(initialInterval time.Duration) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.InitialInterval = initialInterval
	}
}
func WithMultiplier(multiplier float64) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.Multiplier = multiplier
	}
}
func WithMaxInterval(maxInterval time.Duration) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.MaxInterval = maxInterval
	}
}
func WithMaxElapsedTime(maxElapsedTime time.Duration) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.MaxElapsedTime = maxElapsedTime
	}
}
func WithRandomisationFactory(randomisationFactor float64) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.RandomizationFactor = randomisationFactor
	}
}

func WithClock(clock backoff.Clock) Option {
	return func(eb *ExponentialBackoff) {
		eb.expBackoff.Clock = clock
	}
}

func (eb *ExponentialBackoff) Next() (time.Duration, bool) {
	d := eb.expBackoff.NextBackOff()
	if d == backoff.Stop {
		return 0, false
	}

	return d, true
}
func (eb *ExponentialBackoff) After(d time.Duration) <-chan time.Time {
	return eb.afterFunc(d)
}
