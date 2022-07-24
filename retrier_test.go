package recovererr

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type action struct {
	callCounter int
	errors      []error
}

func (a *action) Call() error {
	a.callCounter++
	var err error
	if len(a.errors) > 0 {
		err = a.errors[0]
	}
	if len(a.errors) > 1 {
		a.errors = a.errors[1:]
	}
	return err
}

func Test_retry_retryRecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		retry(ctx, recoverErrorAction.Call, time.Tick(time.Millisecond), RetryRecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("no retry for any other error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		retry(ctx, anyErrorAction.Call, time.Tick(time.Millisecond), RetryRecoverablePolicy)

		assert.Equal(t, 1, anyErrorAction.callCounter)
	})

	t.Run("no retry for unrecoverable error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		retry(ctx, unrecoverableErrorAction.Call, time.Tick(time.Millisecond), RetryRecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		noErrorAction := &action{errors: []error{}}

		retry(ctx, noErrorAction.Call, time.Tick(time.Millisecond), RetryRecoverablePolicy)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
}

func Test_retry_skipUnrecoverablePolicy(t *testing.T) {
	t.Parallel()

	t.Run("retry recoverable until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		recoverErrorAction := &action{errors: []error{Recoverable(errors.New("failure"))}}

		retry(ctx, recoverErrorAction.Call, time.Tick(time.Millisecond), SkipUnrecoverablePolicy)

		assert.True(t, recoverErrorAction.callCounter > 1)
	})

	t.Run("retry any other error until cancellation", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		anyErrorAction := &action{errors: []error{errors.New("failure")}}

		retry(ctx, anyErrorAction.Call, time.Tick(time.Millisecond), SkipUnrecoverablePolicy)

		assert.True(t, anyErrorAction.callCounter > 0)
	})

	t.Run("skip retry for unrecoverable error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		unrecoverableErrorAction := &action{errors: []error{Unrecoverable(errors.New("failure"))}}

		retry(ctx, unrecoverableErrorAction.Call, time.Tick(time.Millisecond), SkipUnrecoverablePolicy)

		assert.Equal(t, 1, unrecoverableErrorAction.callCounter)
	})

	t.Run("no retry on no error", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelFunc()

		noErrorAction := &action{errors: []error{}}

		retry(ctx, noErrorAction.Call, time.Tick(time.Millisecond), SkipUnrecoverablePolicy)

		assert.Equal(t, 1, noErrorAction.callCounter)
	})
}
