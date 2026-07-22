package store

import "errors"

var (
	ErrQueueExists       = errors.New("queue already exists")
	ErrQueueNotFound     = errors.New("queue not found")
	ErrQueueNotEmpty     = errors.New("queue not empty")
	ErrJobNotFound       = errors.New("job not found")
	ErrLeaseNotFound     = errors.New("lease not found")
	ErrLeaseExpired      = errors.New("lease expired")
	ErrInvalidTransition = errors.New("invalid job state transition")
)
