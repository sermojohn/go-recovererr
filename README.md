# go-recovererr
Provides `recoverError` error type that can wrap errors and hold recovery signal.
Also provides `DoRecover` function to check if recovery should be performed for
a given error.

# naming
The package name was conceived by merging `recover` and `error` and can be pronounced as recoverer.
