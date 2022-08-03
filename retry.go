package recovererr

import (
	"context"
	"fmt"
	"time"
)

// DoAndRetry will run a funtion and initiate retries if it fails.
//
// The call to `Retry` is postponed until an error is returned by the function.
func Do(ctx context.Context, f func() error, newBackoffStrategy func() BackoffStrategy, retryPolicy RetryPolicy) error {
	return do(ctx, f, &SystemClock{}, newBackoffStrategy, retryPolicy)
}

func do(ctx context.Context, f func() error, clock Clock, newBackoffStrategy func() BackoffStrategy, retryPolicy RetryPolicy) error {
	err := f()
	if err == nil {
		return nil
	}

	// exit if should not retry
	if !retryPolicy(err) {
		return err
	}

	// initiate backoff strategy
	backoffStrategy := newBackoffStrategy()

	delay, doRetry := backoffStrategy.Next()
	// exit if delay is over
	if !doRetry {
		return err
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("%s, %w", ctx.Err(), err)
	case <-clock.After(delay):
	}

	return retry(ctx, f, clock, backoffStrategy, retryPolicy)
}

// Retry will run the provided function.
//
// If the function fails, retryPolicy is used to extract the recovery context.
// Retry will be performed on intervals provided by a time channel until the context is cancelled.
func Retry(ctx context.Context, f func() error, backoffStrategy BackoffStrategy, retryPolicy RetryPolicy) error {
	return retry(ctx, f, &SystemClock{}, backoffStrategy, retryPolicy)
}

func retry(ctx context.Context, f func() error, clock Clock, backoffStrategy BackoffStrategy, retryPolicy RetryPolicy) error {
	for {
		err := f()
		if err == nil {
			return nil
		}
		// exit if should not retry
		if !retryPolicy(err) {
			return err
		}

		delay, doRetry := backoffStrategy.Next()
		// exit if delay is over
		if !doRetry {
			return err
		}

		// wait or cancel
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s, %w", ctx.Err(), err)
		case <-clock.After(delay):
		}
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

// BackoffStrategy provides different backoff methods for the retry mechanism.
type BackoffStrategy interface {
	Next() (time.Duration, bool)
}

// Clock replaces time package to provide mock replacements.
type Clock interface {
	After(time.Duration) <-chan time.Time
}
