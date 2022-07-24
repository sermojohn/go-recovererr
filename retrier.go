package recovererr

import (
	"context"
	"time"
)

// func retry(ctx context.Context, f func() (bool, error), intervals <-chan time.Time) error {
// 	for {
// 		retry, err := f()
// 		if err == nil {
// 			return nil
// 		}
// 		if !retry {
// 			return err
// 		}

// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case <-intervals:
// 		}
// 		// fmt.Println("will retry")
// 	}
// }

func retry(ctx context.Context, f func() error, intervals <-chan time.Time, retryPolicy RetryPolicy) error {
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

type RetryPolicy func(error) bool

// RetryRecoverablePolicy will return retry if error is explicitly recoverable.
// Any error with no recover context or non matching error will not be retried.
func RetryRecoverablePolicy(err error) bool {
	return DoRecover(err, false)
}

// SkipUnrecoverablePolicy will return not retry if error is explicitly unrecoverable.
// Any error with no recover context will be retried.
func SkipUnrecoverablePolicy(err error) bool {
	return DoRecover(err, true)
}
