package recovererr

import (
	"context"
	"errors"
	"time"
)

// Retry recoverable errors using constant backoff and fixed attempts
func ExampleRetry_First() {
	recoverErrorAction := &mockAction{errors: []error{Recoverable(errors.New("failure"))}}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()
	_ = Retry(ctx, recoverErrorAction.Call, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(5)), RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}

// Retry recoverable errors using exponential backoff.
func ExampleRetry_Second() {
	recoverErrorAction := &mockAction{errors: []error{Recoverable(errors.New("failure"))}}

	mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

	backoffStrategy := NewExponentialBackoff(
		WithInitialInterval(time.Millisecond),
		WithRandomisationFactory(0),
		WithMultiplier(1),
		WithMaxElapsedTime(4*time.Millisecond),
		WithMaxInterval(time.Millisecond),
		WithClock(&mockClock),
	)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancelFunc()
	_ = retry(ctx, recoverErrorAction.Call, &mockClock, backoffStrategy, RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
}

// Retry recoverable errors using constant backoff and default attempts.
func ExampleRetry_Third() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancelFunc()

	recoverErrorAction := &mockAction{errors: []error{Recoverable(errors.New("failure"))}}

	_ = Retry(ctx, recoverErrorAction.Call, NewConstantBackoff(WithInterval(time.Millisecond)), RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
}

// Do run function returning unrecoverable errors using exponential backoff.
func ExampleRetry_Fourth() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	var (
		mockClock                   = mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}
		nonUnrecoverableErrorAction = &mockAction{errors: []error{errors.New("failure")}}
		backoffStretegy             = func() BackoffStrategy {
			return NewExponentialBackoff(
				WithInitialInterval(time.Millisecond),
				WithRandomisationFactory(0),
				WithMultiplier(1),
				WithMaxElapsedTime(2*time.Millisecond),
				WithMaxInterval(time.Millisecond),
				WithClock(&mockClock),
			)
		}
	)

	_ = do(ctx, nonUnrecoverableErrorAction.Call, &mockClock, backoffStretegy, RetryNonUnrecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
}

// Do run function returning unrecoverable errors using constant backoff.
func ExampleRetry_Fifth() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	var (
		nonUnrecoverableErrorAction = &mockAction{errors: []error{errors.New("failure")}}
		backoffStretegy             = func() BackoffStrategy {
			return NewConstantBackoff(
				WithInterval(time.Millisecond),
				WithMaxAttempts(3),
			)
		}
	)

	_ = Do(ctx, nonUnrecoverableErrorAction.Call, backoffStretegy, RetryNonUnrecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
}
