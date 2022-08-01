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

	t.Run("retry recoverable until maximum", func(t *testing.T) {
		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), recoverErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter == 11)
	})

	t.Run("no retry for any other error", func(t *testing.T) {
		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), anyErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.Equal(t, 1, anyErrorAction.callCounter)
	})

	t.Run("no retry for unrecoverable error", func(t *testing.T) {
		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), unrecoverableErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		noErrorAction := &action{errors: []error{}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), noErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryRecoverablePolicy)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
}

func TestRetry_retryNonUnrecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), recoverErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("retry any other error until cancellation", func(t *testing.T) {
		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), anyErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

		assert.True(t, anyErrorAction.callCounter > 0)
	})

	t.Run("skip retry for unrecoverable error", func(t *testing.T) {
		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), unrecoverableErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		noErrorAction := &action{errors: []error{}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), noErrorAction.Call, &mockClock, NewConstantBackoff(time.Millisecond, 10), RetryNonUnrecoverablePolicy)

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
	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	_ = Retry(context.Background(), recoverErrorAction.Call, NewConstantBackoff(time.Millisecond, 5), RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}

func ExampleRetry_Second() {
	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

	backoffStrategy := NewExponentialBackoff(
		WithInitialInterval(time.Millisecond),
		WithRandomisationFactory(0),
		WithMultiplier(1),
		WithMaxElapsedTime(4*time.Millisecond),
		WithMaxInterval(time.Millisecond),
		WithClock(&mockClock),
	)

	_ = retry(context.Background(), recoverErrorAction.Call, &mockClock, backoffStrategy, RetryRecoverablePolicy)

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

func (mc *mockClock) After(time.Duration) <-chan time.Time {
	ch := make(chan time.Time)
	close(ch)
	return ch
}

func ExampleRetry_Third() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancelFunc()

	recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

	_ = Retry(ctx, recoverErrorAction.Call, NewConstantBackoff(time.Millisecond, 5), RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
}

type customError struct {
	recoverable bool
	message     string
}

func (ce *customError) Recover() bool {
	return ce.recoverable
}
func (ce *customError) Error() string {
	return fmt.Sprintf("recoverable:%t, message:%s", ce.recoverable, ce.message)
}

func Test_retry_expired(t *testing.T) {
	t.Parallel()

	t.Run("expired context of recoverable error", func(t *testing.T) {
		action := func() error {
			return &customError{recoverable: true, message: "custom message"}
		}

		backoff := NewConstantBackoff(10*time.Millisecond, 0)
		ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancelFunc()

		err := Retry(ctx, action, backoff, RetryRecoverablePolicy)

		_, recoverable := DoRecover(err)
		assert.True(t, recoverable)
		assert.Equal(t, "context deadline exceeded, recoverable:true, message:custom message", err.Error())

	})
}
