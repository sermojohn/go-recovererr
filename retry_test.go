package recovererr

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry_retryRecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		Retry(ctx, recoverErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("no retry for any other error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		Retry(ctx, anyErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.Equal(t, 1, anyErrorAction.callCounter)
	})

	t.Run("no retry for unrecoverable error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		Retry(ctx, unrecoverableErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		noErrorAction := &action{errors: []error{}}

		Retry(ctx, noErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
}

func TestRetry_retryNonUnrecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		Retry(ctx, recoverErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("retry any other error until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		Retry(ctx, anyErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

		assert.True(t, anyErrorAction.callCounter > 0)
	})

	t.Run("skip retry for unrecoverable error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		Retry(ctx, unrecoverableErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()
		noErrorAction := &action{errors: []error{}}

		Retry(ctx, noErrorAction.Call, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
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

	_ = Retry(ctx, recoverErrorAction.Call, NewConstantBackoff(time.Millisecond, 5), RetryRecoverablePolicy)

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

	mockClosk := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

	backoffStrategy := NewExponentialBackoff(
		WithInitialInterval(time.Millisecond),
		WithRandomisationFactory(0),
		WithMultiplier(1),
		WithMaxElapsedTime(4*time.Millisecond),
		WithMaxInterval(time.Millisecond),
		WithClock(&mockClosk),
	)

	_ = Retry(ctx, recoverErrorAction.Call, backoffStrategy, RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
}

type mockClock struct {
	init     time.Time
	interval time.Duration
}

func (mc *mockClock) Now() time.Time {
	n := mc.init
	mc.init = mc.init.Add(mc.interval)
	return n
}
