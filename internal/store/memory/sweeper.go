package memory

import (
	"context"
	"time"
)

// Sweepers - Will decide later
func (m *Memory) PromoteScheduledJobs(ctx context.Context) (count uint64, err error) {
	return count, nil
}

func (m *Memory) ReapExpiredLeases(ctx context.Context) (count uint64, err error) {
	return count, nil
}

func (m *Memory) DeleteCompletedJobsBefore(ctx context.Context, before time.Time) error {
	return nil
}
