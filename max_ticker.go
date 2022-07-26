package recovererr

import "time"

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
