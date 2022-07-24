package recovererr

import (
	"context"
	"time"
)

// Retry will run the provided function.
// If the function fails, retryPolicy is used to extract
// from the recovery context.
// Retry will be performed on intervals provided by a time channel
// until the context is cancelled.
func Retry(ctx context.Context, f func() error, intervals <-chan time.Time, retryPolicy RetryPolicy) error {
	for {
		err := f()
		if err == nil {
			return nil
		}
		// exit if retry signals not retry
		if !retryPolicy(err) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-intervals:
		}
		// fmt.Println("will retry")
	}
}

// RetryPolicy function implements the policy for performing a retry.
type RetryPolicy func(error) bool

// RetryRecoverablePolicy will return retry if error is recoverable
// and not retry otherwise or for errors with no recovery context.
func RetryRecoverablePolicy(err error) bool {
	found, recover := DoRecover(err)
	return found && recover
}

// RetryNonUnrecoverablePolicy will return retry if error is recoverable
// or error with no recovery context is provided.
func RetryNonUnrecoverablePolicy(err error) bool {
	found, recover := DoRecover(err)
	return !found || recover
}
