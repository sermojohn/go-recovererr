package recovererr

import "strings"

type recoverError struct {
	recover bool
	err     error
}

// Error returns the error in string format.
func (re recoverError) Error() string {
	sb := strings.Builder{}
	if re.recover {
		sb.WriteString("recover: ")
	} else {
		sb.WriteString("unrecover: ")
	}
	sb.WriteString(re.err.Error())
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
