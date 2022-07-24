// Package recoverror provides an error type to signal
// recoverability of an error.
// Error chain is maintained by wrapping and unwrapping
// a verbatim error.
package recoverr

import (
	"errors"
	"strings"
)

type recoverError struct {
	recover bool
	err     error
}

// Error return the error in string format.
func (re recoverError) Error() string {
	sb := &strings.Builder{}
	if re.recover {
		sb.Write([]byte("recover: "))
	} else {
		sb.Write([]byte("unrecover: "))
	}
	sb.Write([]byte(re.err.Error()))
	return sb.String()
}

// Recover provides if should recover from error.
func (re recoverError) Recover() bool {
	return re.recover
}

// Unwrap provides the wrapped error.
func (re recoverError) Unwrap() error {
	return re.err
}

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

// DoRecover is used by error checkers to validate if error should be recovered.
//
// If client marks errors as recoverable using `Recoverable(error)`, this function
// should be used with fallback equal `false` because errors are unrecoverable by default
// and recoverable are explicitly defined.
// If client marks errors as unrecoverable using `Unrecoverable(error)`, this function
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
