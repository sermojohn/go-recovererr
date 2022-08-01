package recovererr

import "time"

// SystemClock provides time package dependency.
type SystemClock struct{}

// After implement the Clock.After method.
func (sc *SystemClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
