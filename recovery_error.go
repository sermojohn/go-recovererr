package recovererr

import "strings"

type recoveryError struct {
	recover bool
	err     error
}

// Error returns the error in string format.
func (re recoveryError) Error() string {
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
func (re recoveryError) Recover() bool {
	return re.recover
}

// Unwrap provides the wrapped error.
func (re recoveryError) Unwrap() error {
	return re.err
}
