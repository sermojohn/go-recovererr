// This package provides utilities to add recovery context to errors.
// Error chain is traversable by unwrapping the given error.
package recovererr

import (
	"errors"
)

// Recoverable wraps an error as recoverable.
func Recoverable(err error) error {
	if err == nil {
		panic("recoverror: error cannot be nil")
	}
	return &recoveryError{err: err, recover: true}
}

// Unrecoverable wraps an error as unrecoverable.
func Unrecoverable(err error) error {
	if err == nil {
		panic("recoverror: error cannot be nil")
	}
	return &recoveryError{err: err}
}

// DoRecover can be used by user to validate if error should be recovered.
// When no recovery context is found in the given error, it returns
// false in both values.
// When recovery context is found in the given error, it returns
// true in the first value and the recovery context of the error.
func DoRecover(err error) (bool, bool) {
	for err != nil {
		if x, ok := err.(interface{ Recover() bool }); ok {
			return true, x.Recover()
		}
		err = errors.Unwrap(err)
	}
	return false, false
}
