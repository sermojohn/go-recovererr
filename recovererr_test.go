package recovererr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoRecover(t *testing.T) {
	var (
		ConnectionError = errors.New("connection error")
		ParseError      = errors.New("parse error")
	)
	t.Run("recoverable wrapped error", func(t *testing.T) {
		err := Recoverable(ConnectionError)
		found, recover := DoRecover(err)
		assert.True(t, found, err)
		assert.True(t, recover, err)
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		err := Unrecoverable(ParseError)
		found, recover := DoRecover(err)
		assert.True(t, found, err)
		assert.False(t, recover, err)
	})

	t.Run("recoverable wrapped error wrapped", func(t *testing.T) {
		err := fmt.Errorf("failed to store object, %w", Recoverable(ConnectionError))
		found, recover := DoRecover(err)
		assert.True(t, found, err)
		assert.True(t, recover, err)
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		err := fmt.Errorf("failed to parse object, %w", Unrecoverable(ParseError))
		found, recover := DoRecover(err)
		assert.True(t, found, err)
		assert.False(t, recover, err)
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		err := &anyError{}
		found, recover := DoRecover(err)
		assert.False(t, found, err)
		assert.False(t, recover, err)
	})

	t.Run("other recover error implementation", func(t *testing.T) {
		err := &otherRecoverError{recover: true}
		found, recover := DoRecover(err)
		assert.True(t, found, err)
		assert.True(t, recover, err)
	})
}

func TestRecoverable(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover(), "expected pacic")
		}()

		Recoverable(nil)
	})

	t.Run("not-nil error", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover(), "expected no panic")
		}()

		Recoverable(errors.New("not-nil error"))
	})
}

func TestUnrecoverable(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		defer func() {
			assert.NotNil(t, recover(), "expected pacic")
		}()

		Unrecoverable(nil)
	})

	t.Run("not-nil error", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover(), "expected no panic")
		}()

		Unrecoverable(errors.New("not-nil error"))
	})
}

type anyError struct{}

func (ae anyError) Error() string { return "" }

type otherRecoverError struct {
	recover bool
}

func (ore otherRecoverError) Error() string {
	return fmt.Sprintf("recover: %t", ore.recover)
}
func (ore otherRecoverError) Recover() bool {
	return ore.recover
}
