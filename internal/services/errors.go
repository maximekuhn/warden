package services

import "errors"

// Non-exhaustive list of errors returned by different services.
var (
	ErrNoPortAvailable = errors.New("no port available")
)
