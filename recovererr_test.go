package recoverr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ConnectionError = errors.New("connection error")
	ParseError      = errors.New("parse error")
)

func TestDoRecover(t *testing.T) {
	t.Run("recoverable wrapped error", func(t *testing.T) {
		err := Recoverable(ConnectionError)
		assert.True(t, DoRecover(err, false), err)
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		err := Unrecoverable(ParseError)
		assert.False(t, DoRecover(err, true), err)
	})

	t.Run("recoverable wrapped error wrapped", func(t *testing.T) {
		err := fmt.Errorf("failed to store object, %w", Recoverable(ConnectionError))
		assert.True(t, DoRecover(err, false), err)
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		err := fmt.Errorf("failed to parse object, %w", Unrecoverable(ParseError))
		assert.False(t, DoRecover(err, true), err)
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		err := &anyError{}
		assert.False(t, DoRecover(err, false), err)
	})

	t.Run("other recover error implementation", func(t *testing.T) {
		err := &otherRecoverError{recover: true}
		assert.True(t, DoRecover(err, false), err)
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("recoverable connection error", func(t *testing.T) {
		err := Recoverable(ConnectionError)
		assert.Equal(t, ConnectionError, errors.Unwrap(err))
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		err := Unrecoverable(ParseError)
		assert.Equal(t, ParseError, errors.Unwrap(err))
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
