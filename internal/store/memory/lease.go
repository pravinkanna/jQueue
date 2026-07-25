package memory

import (
	"context"
	"time"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Lease
func (m *Memory) LeaseJobs(ctx context.Context, queue string, batchSize uint32, leaseDuration time.Duration) (leasedJobs []store.LeasedJob, err error) {
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
