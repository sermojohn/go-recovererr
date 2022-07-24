# go-recovererr
Provides `Recoverable(error) error` and `Unrecoverable(error) error` functions that wrap the given error with recovery context.

Also provides `DoRecover` function to check the recovery context of any error.

# naming
The package name was conceived by merging `recover` and `error` and can be pronounced as recoverer.
