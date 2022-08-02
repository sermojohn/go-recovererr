package recovererr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("recoverable connection error", func(t *testing.T) {
		var connectionError = errors.New("connection error")
		err := recoverError{recover: true, err: connectionError}
		assert.Equal(t, connectionError, errors.Unwrap(err))
	})

	t.Run("unrecoverable wrapped error", func(t *testing.T) {
		var ParseError = errors.New("parse error")
		err := Unrecoverable(ParseError)
		assert.Equal(t, ParseError, errors.Unwrap(err))
	})
}

func TestRecoverError_Error(t *testing.T) {
	t.Parallel()

	t.Run("recoverable error", func(t *testing.T) {
		re := &recoverError{recover: true, err: errors.New("test")}

		assert.Equal(t, "recover: test", re.Error())
	})

	t.Run("unrecoverable error", func(t *testing.T) {
		re := &recoverError{recover: false, err: errors.New("test")}

		assert.Equal(t, "unrecover: test", re.Error())
	})
}
