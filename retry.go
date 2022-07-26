package recovererr

import (
	"context"
	"time"
)

// NewRetrier creates a new retrier and provides dynamic options.
//
// If required options are empty, defaults apply.
func NewRetrier(opts ...Option) *Retrier {
	r := Retrier{}

	for _, opt := range append(opts, withDefaults()) {
		opt(&r)
	}

	return &r
}

// Retrier holds the static parameters of a retry mechanism
type Retrier struct {
	intervalGenerator IntervalGenerator
	retryPolicy       RetryPolicy
}

// Option to customise retrier.
type Option func(*Retrier)

// WithRetryPolicy applies the provided retry policy to the retrier.
func WithRetryPolicy(retryPolicy RetryPolicy) Option {
	return func(r *Retrier) {
		r.retryPolicy = retryPolicy
	}
}

// WithIntervalGenerator applies the provided interval generator to the retrier.
func WithIntervalGenerator(intervalGenerator IntervalGenerator) Option {
	return func(r *Retrier) {
		r.intervalGenerator = intervalGenerator
	}
}

func withDefaults() Option {
	return func(r *Retrier) {
		if r.retryPolicy == nil {
			r.retryPolicy = RetryRecoverablePolicy
		}
		if r.intervalGenerator == nil {
			r.intervalGenerator = StaticIntervalGenerator(1)
		}
	}
}

// Do will run the provided function.
// If the function fails, retryPolicy is used to extract
// from the recovery context.
// Retry will be performed on intervals provided by a time channel
// until the context is cancelled.
func (r *Retrier) Do(ctx context.Context, f func() error) error {
	intervalInitFunc, intervalCancelFunc := r.intervalGenerator()
	defer intervalCancelFunc()

	var intervalChan <-chan time.Time
	zeroTime := time.Time{}
	for {
		// call func
		err := f()
		if err == nil {
			return nil
		}

		// exit if retry signals not retry
		if !r.retryPolicy(err) {
			return err
		}

		// initialise the interval generator
		if intervalChan == nil {
			intervalChan = intervalInitFunc()
		}

		// timeout or wait
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-intervalChan:
			if t == zeroTime {
				return nil
			}
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
