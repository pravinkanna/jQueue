package memory

import (
	"context"
	"time"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Lease
func (m *Memory) LeaseJob(ctx context.Context, queue string, leaseDuration time.Duration) (leasedJobs store.LeasedJob, err error) {
	// check the queue exist

	// make sure the queue len
	return leasedJobs, nil
}

func (m *Memory) ExtendJobLease(ctx context.Context, leaseToken string, duration time.Duration) (leaseExpiresAt time.Time, err error) {
	return leaseExpiresAt, nil
}

func (m *Memory) AckJob(ctx context.Context, leaseToken string) error {
	return nil
}

func (m *Memory) NackJob(ctx context.Context, leaseToken string, reason string) error {
	return nil
}
