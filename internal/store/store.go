package store

import (
	"context"
	"time"
)

type Queue struct {
	Name           string
	PendingCount   uint64
	ScheduledCount uint64
	LeasedCount    uint64
	CompletedCount uint64
	FailedCount    uint64
	DLQCount       uint64
}

type EnqueueParams struct {
	IdempotencyKey string
	Queue          string
	Payload        []byte
	MaxRetries     uint32
	RunAt          time.Time
}

type JobState int

const (
	JobStateUnspecified JobState = iota
	JobStateScheduled
	JobStatePending
	JobStateLeased
	JobStateCompleted
	JobStateFailed
	JobStateDLQ
	JobStateCancelled
)

type Job struct {
	JobID          string
	IdempotencyKey string
	Queue          string
	Payload        []byte
	State          JobState
	MaxRetries     uint32
	AttemptCount   uint32
	LastError      string
	CreatedAt      time.Time
	ScheduledAt    time.Time
	CompletedAt    time.Time
}

type LeasedJob struct {
	LeaseToken string
	Job        Job
	ExpiresAt  time.Time
}

type Store interface {
	// Queue
	CreateQueue(ctx context.Context, name string) error
	DeleteQueue(ctx context.Context, name string) error
	PurgeQueue(ctx context.Context, name string) (purgedCount uint64, err error)
	ListQueues(ctx context.Context) (queues []Queue, err error)
	GetQueueStatus(ctx context.Context, name string) (queue Queue, err error)

	// Job
	EnqueueJob(ctx context.Context, enqueueParams EnqueueParams) (jobID string, isDuplicate bool, err error)
	GetJob(ctx context.Context, jobID string) (job Job, err error)
	CancelJob(ctx context.Context, jobID string) (state JobState, err error)
	ListJobs(ctx context.Context, queue string, state JobState, pageSize uint32, pageToken string) (jobs []Job, nextPageToken string, err error)
	RetryDLQJob(ctx context.Context, jobID string) error
	RetryDLQJobs(ctx context.Context, queue string) (retriedCount uint32, err error)

	// Lease
	LeaseJobs(ctx context.Context, queue string, batchSize uint32, leaseDuration time.Duration) (leasedJobs []LeasedJob, err error)
	ExtendJobLease(ctx context.Context, leaseToken string, duration time.Duration) (leaseExpiresAt time.Time, err error)
	AckJob(ctx context.Context, leaseToken string) error
	NackJob(ctx context.Context, leaseToken string, reason string) (job Job, err error)

	// Sweepers - Will decide later
	PromoteScheduledJobs(ctx context.Context) (count uint64, err error)
	ReapExpiredLeases(ctx context.Context) (count uint64, err error)
	DeleteCompletedJobsBefore(ctx context.Context, before time.Time) error
}
