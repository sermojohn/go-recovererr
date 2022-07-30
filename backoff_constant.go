package recovererr

import "time"

type ConstantBackoff struct {
	delay     time.Duration
	afterFunc func(time.Duration) <-chan time.Time
	max       int

	counter int
}

func NewConstantBackoff(d time.Duration, max int) *ConstantBackoff {
	cb := ConstantBackoff{
		delay:     d,
		afterFunc: time.After,
		max:       max,
	}
	return &cb
}

func (cb *ConstantBackoff) Next() (time.Duration, bool) {
	cb.counter++
	if cb.counter > cb.max {
		return 0, false
	}
	return cb.delay, true
}

func (cb *ConstantBackoff) After(d time.Duration) <-chan time.Time {
	return cb.afterFunc(d)
}
