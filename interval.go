package recovererr

import (
	"time"

	"github.com/cenkalti/backoff"
)

// IntervalGenerator provides a time channel with intervals and a cancel function.
type IntervalGenerator func() (func() <-chan time.Time, func())

func BackoffGenerator() (func() <-chan time.Time, func()) {
	backoffer := &backoff.ExponentialBackOff{
		InitialInterval:     1 * time.Millisecond,
		RandomizationFactor: 0,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         10 * time.Millisecond,
		MaxElapsedTime:      10 * time.Millisecond,
		Clock:               backoff.SystemClock,
	}
	ticker := backoff.NewTicker(backoffer)

	return func() <-chan time.Time {
			backoffer.Reset()
			return ticker.C
		},
		func() { ticker.Stop() }
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

type MaxTicker struct {
	max              int
	interval         time.Duration
	ticker           *time.Ticker
	tickerCancelFunc func()
	c                chan time.Time
}

func NewMaxTicker(max int, d time.Duration) *MaxTicker {
	ticker := time.NewTicker(d)
	mt := &MaxTicker{
		max:    max,
		ticker: ticker,
		tickerCancelFunc: func() {
			ticker.Stop()
		},
		c: make(chan time.Time, 1),
	}
	go mt.run()
	return mt
}

func (mt *MaxTicker) run() {
	var counter int
	for {
		select {
		case t := <-mt.ticker.C:
			// Give the channel a 1-element time buffer.
			// If the client falls behind while reading, we drop ticks
			// on the floor until the client catches up.
			select {
			case mt.c <- t:
				counter++
				if counter == mt.max {
					close(mt.c)
					return
				}
			default:
			}
		}
	}
}

func (mt *MaxTicker) Stop() {
	mt.tickerCancelFunc()
}
