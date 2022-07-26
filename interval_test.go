package recovererr

import (
	"context"
	"errors"
	"time"
)

func ExampleRetry_Second() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()

	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	r := NewRetrier(
		WithRetryPolicy(RetryRecoverablePolicy),
		WithIntervalGenerator(BackoffGenerator(1*time.Millisecond, 0, 1.5, 10*time.Millisecond, 10*time.Millisecond)))
	r.Do(ctx, recoverErrorAction.Call)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}

func ExampleRetry_First() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()

	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	r := NewRetrier(
		WithRetryPolicy(RetryRecoverablePolicy),
		WithIntervalGenerator(StaticAttemptsGenerator(5, 1*time.Nanosecond)))
	r.Do(ctx, recoverErrorAction.Call)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}

func ExampleRetry_Third() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()

	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	r := NewRetrier(
		WithRetryPolicy(RetryRecoverablePolicy),
		WithIntervalGenerator(StaticIntervalGenerator(4*time.Millisecond)))
	r.Do(ctx, recoverErrorAction.Call)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
}
