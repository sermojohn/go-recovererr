package recovererr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry_retryRecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until maximum", func(t *testing.T) {
		recoverErrorAction := &mockAction{errors: []error{Recoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), recoverErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryRecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter == 11)
	})

	t.Run("no retry for any other error", func(t *testing.T) {
		anyErrorAction := &mockAction{errors: []error{errors.New("failure")}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), anyErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryRecoverablePolicy)

		assert.Equal(t, 1, anyErrorAction.callCounter)
	})

	t.Run("no retry for unrecoverable error", func(t *testing.T) {
		unrecoverableErrorAction := &mockAction{errors: []error{Unrecoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), unrecoverableErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryRecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		noErrorAction := &mockAction{errors: []error{}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), noErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryRecoverablePolicy)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})

	t.Run("retry forever on no error", func(t *testing.T) {
		noErrorAction := &mockAction{errors: []error{}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), noErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(9)), RetryForever)

		assert.Equal(t, 10, noErrorAction.callCounter)
	})
}

func TestRetry_retryNonUnrecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		recoverErrorAction := &mockAction{errors: []error{Recoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), recoverErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryNonUnrecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("retry any other error until cancellation", func(t *testing.T) {
		anyErrorAction := &mockAction{errors: []error{errors.New("failure")}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), anyErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryNonUnrecoverablePolicy)

		assert.True(t, anyErrorAction.callCounter > 0)
	})

	t.Run("skip retry for unrecoverable error", func(t *testing.T) {
		unrecoverableErrorAction := &mockAction{errors: []error{Unrecoverable(errors.New("failure"))}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), unrecoverableErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryNonUnrecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		noErrorAction := &mockAction{errors: []error{}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		_ = retry(context.Background(), noErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10)), RetryNonUnrecoverablePolicy)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
}

type mockAction struct {
	callCounter int
	errors      []error
}

func (ma *mockAction) Call() error {
	ma.callCounter++
	fmt.Printf("action called %d time(s)\n", ma.callCounter)

	var err error
	if len(ma.errors) > 0 {
		err = ma.errors[0]
	}
	if len(ma.errors) > 1 {
		ma.errors = ma.errors[1:]
	}
	return err
}

func ExampleRetry_First() {
	recoverErrorAction := &mockAction{errors: []error{Recoverable(errors.New("failure"))}}

	_ = Retry(context.Background(), recoverErrorAction.Call, NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(5)), RetryRecoverablePolicy)

	// Output:
	// action called 1 time(s)
	// action called 2 time(s)
	// action called 3 time(s)
	// action called 4 time(s)
	// action called 5 time(s)
	// action called 6 time(s)
}

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

	recoverErrorAction := &mockAction{errors: []error{Recoverable(errors.New("failure"))}}

	_ = Retry(ctx, recoverErrorAction.Call, NewConstantBackoff(WithInterval(time.Millisecond)), RetryRecoverablePolicy)

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

func TestRetry_expired(t *testing.T) {
	t.Parallel()

	t.Run("expired context of recoverable error", func(t *testing.T) {
		action := func() error {
			return &customError{recoverable: true, message: "custom message"}
		}

		backoff := NewConstantBackoff()
		ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancelFunc()

		err := Retry(ctx, action, backoff, RetryRecoverablePolicy)

		_, recoverable := DoRecover(err)
		assert.True(t, recoverable)
		assert.Equal(t, "context deadline exceeded, recoverable:true, message:custom message", err.Error())

	})
}

func TestDo_expired(t *testing.T) {
	t.Parallel()

	t.Run("expired context of recoverable error", func(t *testing.T) {
		action := func() error {
			return &customError{recoverable: true, message: "custom message"}
		}

		newBackoff := func() BackoffStrategy {
			cb := NewConstantBackoff()
			return cb
		}
		ctx, cancelFunc := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancelFunc()

		err := Do(ctx, action, newBackoff, RetryRecoverablePolicy)

		_, recoverable := DoRecover(err)
		assert.True(t, recoverable)
		assert.Equal(t, "context deadline exceeded, recoverable:true, message:custom message", err.Error())

	})
}

func Test_do(t *testing.T) {
	type args struct {
		ctx                 context.Context
		f                   *mockAction
		clock               Clock
		backoffStrategyFunc func() BackoffStrategy
		retryPolicy         RetryPolicy
	}
	tests := []struct {
		name       string
		args       args
		assertFunc func(*testing.T, *args, error)
	}{
		{
			name: "retry recoverable policy no retry on no error",
			args: args{
				f:           &mockAction{},
				retryPolicy: RetryRecoverablePolicy,
			},
			assertFunc: func(t *testing.T, a *args, err error) {},
		},
		{
			name: "retry recoverable policy return error error",
			args: args{
				f:           &mockAction{errors: []error{Unrecoverable(fmt.Errorf("test"))}},
				retryPolicy: RetryRecoverablePolicy,
			},
			assertFunc: func(t *testing.T, a *args, err error) {
				assert.Equal(t, Unrecoverable(fmt.Errorf("test")), err)
			},
		},
		{
			name: "retry recoverable policy on recoverable error",
			args: args{
				ctx:         context.Background(),
				f:           &mockAction{errors: []error{Recoverable(fmt.Errorf("test"))}},
				retryPolicy: RetryRecoverablePolicy,
				backoffStrategyFunc: func() BackoffStrategy {
					return NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(3))
				},
				clock: &mockClock{init: time.Now(), interval: time.Nanosecond},
			},
			assertFunc: func(t *testing.T, a *args, err error) {
				assert.Equal(t, Recoverable(fmt.Errorf("test")), err)
			},
		},
		{
			name: "rretry recoverable policy retry until recovered",
			args: args{
				ctx:         context.Background(),
				f:           &mockAction{errors: []error{Recoverable(fmt.Errorf("test")), nil}},
				retryPolicy: RetryRecoverablePolicy,
				backoffStrategyFunc: func() BackoffStrategy {
					return NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(10))
				},
				clock: &mockClock{init: time.Now(), interval: time.Nanosecond},
			},
			assertFunc: func(t *testing.T, a *args, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, a.f.callCounter)
			},
		},
		{
			name: "retry recoverable policy retry until expiration",
			args: args{
				ctx:         contextWithTimeout(time.Millisecond),
				f:           &mockAction{errors: []error{Recoverable(io.EOF)}},
				retryPolicy: RetryRecoverablePolicy,
				backoffStrategyFunc: func() BackoffStrategy {
					return NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(3))
				},
				clock: &SystemClock{},
			},
			assertFunc: func(t *testing.T, a *args, err error) {
				assert.True(t, errors.Is(err, io.EOF), err)

			},
		},
		{
			name: "retry recoverable policy no retry on finished backoff",
			args: args{
				ctx:         contextWithTimeout(time.Millisecond),
				f:           &mockAction{errors: []error{Recoverable(io.EOF)}},
				retryPolicy: RetryRecoverablePolicy,
				backoffStrategyFunc: func() BackoffStrategy {
					cb := NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(1))
					// make the backoff complete
					cb.Next()
					return cb
				},
				clock: &SystemClock{},
			},
			assertFunc: func(t *testing.T, a *args, err error) {
				assert.True(t, errors.Is(err, io.EOF), err)
				assert.Equal(t, 1, a.f.callCounter)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := do(tt.args.ctx, tt.args.f.Call, tt.args.clock, tt.args.backoffStrategyFunc, tt.args.retryPolicy)

			tt.assertFunc(t, &tt.args, err)
		})
	}
}

func contextWithTimeout(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
