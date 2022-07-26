package recovererr

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetrier_Do_RetryRecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		retrier := NewRetrier(WithRetryPolicy(RetryRecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, recoverErrorAction.Call)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("no retry for any other error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		retrier := NewRetrier(WithRetryPolicy(RetryRecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, anyErrorAction.Call)

		assert.Equal(t, 1, anyErrorAction.callCounter)
	})

	t.Run("no retry for unrecoverable error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		retrier := NewRetrier(WithRetryPolicy(RetryRecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, unrecoverableErrorAction.Call)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		noErrorAction := &action{errors: []error{}}

		retrier := NewRetrier(WithRetryPolicy(RetryRecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, noErrorAction.Call)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
}

func TestRetier_Do_RetryNonUnrecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		retrier := NewRetrier(WithRetryPolicy(RetryNonUnrecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, recoverErrorAction.Call)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("retry any other error until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		retrier := NewRetrier(WithRetryPolicy(RetryNonUnrecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, anyErrorAction.Call)

		assert.True(t, anyErrorAction.callCounter > 0)
	})

	t.Run("skip retry for unrecoverable error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		retrier := NewRetrier(WithRetryPolicy(RetryNonUnrecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, unrecoverableErrorAction.Call)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		noErrorAction := &action{errors: []error{}}

		retrier := NewRetrier(WithRetryPolicy(RetryNonUnrecoverablePolicy), WithIntervalGenerator(StaticIntervalGenerator(time.Millisecond)))
		retrier.Do(ctx, noErrorAction.Call)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
}

func TestRetrier_withDefaults(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancelFunc()
	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	retrier := NewRetrier()
	retrier.Do(ctx, recoverErrorAction.Call)

	assert.True(t, recoverErrorAction.callCounter > 1)
}

type action struct {
	callCounter int
	errors      []error
}

func (a *action) Call() error {
	a.callCounter++
	fmt.Printf("action called %d time(s)\n", a.callCounter)

	var err error
	if len(a.errors) > 0 {
		err = a.errors[0]
	}
	if len(a.errors) > 1 {
		a.errors = a.errors[1:]
	}
	return err
}

func ExampleRetry_First() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()

	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	r := NewRetrier(
		WithRetryPolicy(RetryRecoverablePolicy),
		WithIntervalGenerator(StaticAttemptsGenerator(5, 1)))
	r.Do(ctx, recoverErrorAction.Call)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}

func ExampleRetry_Second() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelFunc()

	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	r := NewRetrier(
		WithRetryPolicy(RetryRecoverablePolicy),
		WithIntervalGenerator(BackoffGenerator))
	r.Do(ctx, recoverErrorAction.Call)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}
