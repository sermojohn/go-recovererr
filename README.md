# go-recovererr [![Go Reference](https://pkg.go.dev/badge/github.com/sermojohn/go-recovererr.svg)](https://pkg.go.dev/github.com/sermojohn/go-recovererr) [![Go Report Card](https://goreportcard.com/badge/github.com/sermojohn/go-recovererr)](https://goreportcard.com/report/github.com/sermojohn/go-recovererr)

## Installation
```
go get -u github.com/sermojohn/go-recovererr
```

## Usage
TBD

### Naming
The package name was conceived by merging `recover` and `error` and can be pronounced as recoverer.

## Description

### Error recovery context
The package provides `Recoverable` and `Unrecoverable` public functions to wrap the given error with recovery context.
Also provides `DoRecover` function to check the recovery context of any error.

### Retry
The package provides `Retry` function that receives a function that can return an error. 
The `RetryPolicy` is also provided to check the error recovery context on failure and define if the function should be retried.
`Retry` will perform a retry on function failure, after the intervals channel fires, until the provided context.Context is cancelled.


