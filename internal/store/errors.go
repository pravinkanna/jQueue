package store

import "errors"

var (
	ErrQueueExists                  = errors.New("queue already exists")
	ErrQueueNotFound                = errors.New("queue not found")
	ErrQueueNotEmpty                = errors.New("queue not empty")
	ErrQueueEmpty                   = errors.New("queue empty")
	ErrJobNotFound                  = errors.New("job not found")
	ErrJobNotLeased                 = errors.New("job not leased")
	ErrJobAlreadyStartedOrCompleted = errors.New("job already started or completed")
	ErrLeaseNotFound                = errors.New("lease not found")
	ErrLeaseExpired                 = errors.New("lease expired")
	ErrInvalidTransition            = errors.New("invalid job state transition")
)
