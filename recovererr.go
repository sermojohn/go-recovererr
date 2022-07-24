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
	return &recoverError{err: err, recover: true}
}

// Unrecoverable wraps an error as unrecoverable.
func Unrecoverable(err error) error {
	if err == nil {
		panic("recoverror: error cannot be nil")
	}
	return &recoverError{err: err}
}

// DoRecover can be used by error checkers to validate if error should be recovered.
//
// If the user marks errors as recoverable using `Recoverable(error)`, this function
// should be used with fallback equal `false` because errors are unrecoverable by default
// and recoverable are explicitly defined.
// If the user marks errors as unrecoverable using `Unrecoverable(error)`, this function
// should be used with fallback equal `true` because errors are recoverable by default
// and unrecoverable are explicitly defined.
func DoRecover(err error, fallback bool) bool {
	for err != nil {
		if x, ok := err.(interface{ Recover() bool }); ok {
			return x.Recover()
		}
		err = errors.Unwrap(err)
	}
	return fallback
}
