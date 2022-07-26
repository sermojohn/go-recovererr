package recovererr

import (
	"time"

	"github.com/cenkalti/backoff"
)

// IntervalGenerator provides a time channel with intervals and a cancel function.
type IntervalGenerator func() (func() <-chan time.Time, func())

func BackoffGenerator(initialInterval time.Duration, randomization float64, multiplier float64, maxInterval time.Duration, maxElapsedTime time.Duration) func() (func() <-chan time.Time, func()) {
	backoffer := &backoff.ExponentialBackOff{
		InitialInterval:     initialInterval,
		RandomizationFactor: randomization,
		Multiplier:          multiplier,
		MaxInterval:         maxInterval,
		MaxElapsedTime:      maxElapsedTime,
		Clock:               backoff.SystemClock,
	}
	ticker := backoff.NewTicker(backoffer)

	return func() (func() <-chan time.Time, func()) {
		return func() <-chan time.Time {
				backoffer.Reset()
				return ticker.C
			},
			func() { ticker.Stop() }
	}
}

func StaticIntervalGenerator(d time.Duration) func() (func() <-chan time.Time, func()) {
	return func() (func() <-chan time.Time, func()) {
		ticker := time.NewTicker(d)
		return func() <-chan time.Time {
			return ticker.C
		}, func() { ticker.Stop() }
	}
}

func StaticAttemptsGenerator(max int, d time.Duration) func() (func() <-chan time.Time, func()) {
	return func() (func() <-chan time.Time, func()) {
		ticker := NewMaxTicker(max, d)
		return func() <-chan time.Time { return ticker.c }, func() { ticker.Stop() }
	}
}
