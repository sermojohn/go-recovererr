package recovererr

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff"
)

func ExampleRetry_ExponentialBackoff() {
	backoffer := &backoff.ExponentialBackOff{
		InitialInterval:     1 * time.Millisecond,
		RandomizationFactor: 0,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         10 * time.Millisecond,
		MaxElapsedTime:      10 * time.Millisecond,
		Clock:               backoff.SystemClock,
	}
	backoffer.Reset()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()
	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	Retry(ctx, recoverErrorAction.Call, backoff.NewTicker(backoffer).C, RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}
