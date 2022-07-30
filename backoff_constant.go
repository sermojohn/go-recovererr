package recovererr

import "time"

type ConstantBackoff struct {
	delay     time.Duration
	afterFunc func(time.Duration) <-chan time.Time
	max       int

	counter int
}

func (cb *ConstantBackoff) Next() (time.Duration, bool) {
	cb.counter++
	if cb.counter > cb.max {
		return 0, false
	}
	return cb.delay, true
}

func (cb *ConstantBackoff) After(d time.Duration) <-chan time.Time {
	if cb.counter > cb.max {
		ch := make(chan time.Time)
		close(ch)
		return ch
	}
	return cb.afterFunc(d)
}
