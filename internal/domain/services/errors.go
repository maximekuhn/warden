package services

import "errors"

// Non-exhaustive list of errors returned by different services.
var (
	ErrNoPortAvailable               = errors.New("no port available")
	ErrServerAlreadyHasAllocatedPort = errors.New("server already has an allocated port")
	ErrServerNotFound                = errors.New("minecraft server not found")
)
