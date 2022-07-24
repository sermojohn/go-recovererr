package recovererr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrap(t *testing.T) {
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
