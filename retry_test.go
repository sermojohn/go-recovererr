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

	t.Run("cancelled context after no error", func(t *testing.T) {
		noErrorAction := &mockAction{errors: []error{}}

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		err := retry(&mockContext{done: true}, noErrorAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond)), RetryForever)

		assert.Nil(t, err)
	})

	t.Run("cancelled context after call error", func(t *testing.T) {
		var (
			actionError  = errors.New("failed!")
			failedAction = &mockAction{errors: []error{actionError}}
		)

		mockClock := mockClock{init: time.Unix(1659219915, 0), interval: time.Millisecond}

		err := retry(&mockContext{done: true}, failedAction.Call, &mockClock, NewConstantBackoff(WithInterval(time.Millisecond)), RetryForever)

		assert.True(t, errors.Is(err, actionError))
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
	t.Parallel()

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
		{
			name: "cancel context after successful call",
			args: args{
				ctx:         &mockContext{done: true},
				f:           &mockAction{},
				retryPolicy: RetryRecoverablePolicy,
				backoffStrategyFunc: func() BackoffStrategy {
					return NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(1))
				},
				clock: &SystemClock{},
			},
			assertFunc: func(t *testing.T, a *args, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "cancel context after failed call",
			args: args{
				ctx:         &mockContext{done: true},
				f:           &mockAction{errors: []error{errors.New("call failed")}},
				retryPolicy: RetryRecoverablePolicy,
				backoffStrategyFunc: func() BackoffStrategy {
					return NewConstantBackoff(WithInterval(time.Millisecond), WithMaxAttempts(1))
				},
				clock: &SystemClock{},
			},
			assertFunc: func(t *testing.T, a *args, err error) {
				assert.Equal(t, err, errors.New("call failed"))
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

type mockContext struct {
	done         bool
	err          error
	val          interface{}
	deadlineTime time.Time
	deadlineSet  bool
}

func (mc *mockContext) Deadline() (deadline time.Time, ok bool) {
	return mc.deadlineTime, mc.deadlineSet
}

func (mc *mockContext) Done() <-chan struct{} {
	ch := make(chan struct{})
	if mc.done {
		close(ch)
	}
	return ch
}

func (mc *mockContext) Err() error {
	return mc.err
}

func (mc *mockContext) Value(_ interface{}) interface{} {
	return mc.val
}
