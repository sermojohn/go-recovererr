[![Go Reference](https://pkg.go.dev/badge/github.com/sermojohn/go-recovererr.svg)](https://pkg.go.dev/github.com/sermojohn/go-recovererr)
[![Go Report Card](https://goreportcard.com/badge/github.com/sermojohn/go-recovererr)](https://goreportcard.com/report/github.com/sermojohn/go-recovererr)

# go-recovererr

# naming
The package name was conceived by merging `recover` and `error` and can be pronounced as recoverer.

# recover error
Provides `Recoverable(error) error` and `Unrecoverable(error) error` functions that wrap the given error with recovery context.
Also provides `DoRecover` function to check the recovery context of any error.

# retry
Provides `Retry` function that receives an action that can return an error. Additionally receives a `RetryPolicy` that checks
the recovery context of the error and defines if the action should be retried.
It will retry on failure every time the intervals channel fires, until the provided context.Context is cancelled.
