# go-recovererr [![Go Reference](https://pkg.go.dev/badge/github.com/sermojohn/go-recovererr.svg)](https://pkg.go.dev/github.com/sermojohn/go-recovererr) [![Go Report Card](https://goreportcard.com/badge/github.com/sermojohn/go-recovererr)](https://goreportcard.com/report/github.com/sermojohn/go-recovererr)

## Installation
```
go get -u github.com/sermojohn/go-recovererr
```

## Naming
The package name was conceived by merging `recover` and `error` and can be pronounced as recoverer.


## Usage
1. Wrap recoverable errors with an error value implementing `Recover() bool`:
```
type customError struct {
	recoverable bool
	message     string
}

func (ce *customError) Recover() bool {
	return ce.recoverable
}
func (ce *customError) Error() string {
	return fmt.Sprintf("recoverable:%t, message:%s", ce.recoverable, ce.message)
}
```

2. Retry action of recoverable error per 1sec forever:
```
action := func() error {
    return &customError{recoverable: true}
}

backoff := NewConstantBackoff(time.Second, 0)

_ = Retry(context.Background(), action, backoff, RetryRecoverablePolicy)
```

3. Retry action of recoverable error using exponential backoff starting with 1sec until 5sec elapse:
```
action := func() error {
    return &customError{recoverable: true}
}
backoff := NewExponentialBackoff(
    WithInitialInterval(time.Second), 
    WithMaxElapsedTime(5*time.Second),
)

_ = Retry(context.Background(), action, backoff, RetryRecoverablePolicy)
```

4. Retry action of non unrecoverable error:
```
action := func() error {
    return errors.New("any error")
}
backoff := NewConstantBackoff(time.Second, 0)

_ = Retry(context.Background(), action, backoff, RetryNonUnrecoverablePolicy)
```
## Description

### Error recovery context
The package provides `Recoverable` and `Unrecoverable` public functions to wrap the given error with recovery context.
Also provides `DoRecover` function to check the recovery context of any error.

### Retry
The package provides function `Retry` that receives a function that optionally returns an error. 
The `RetryPolicy` is provided to `Retry`, to check the error recovery context on failure and define if the function should be retried.
The `BackoffStrategy` is provided to defind the delay applied before each retry performing either `constant` or `exponential` backoff.
If context.Context gets cancelled no extra retry will be performed, but the original error will be wrapped to the timeout error.


