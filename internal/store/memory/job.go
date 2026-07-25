package memory

import (
	"context"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Job
func (m *Memory) EnqueueJob(ctx context.Context, enqueueParams store.EnqueueParams) (jobID string, isDuplicate bool, err error) {
	return jobID, isDuplicate, nil
}

func (m *Memory) GetJob(ctx context.Context, jobID string) (job store.Job, err error) {
	return job, nil
}

func (m *Memory) CancelJob(ctx context.Context, jobID string) (state store.JobState, err error) {
	return state, nil
}

func (m *Memory) ListJobs(ctx context.Context, queue string, state store.JobState, pageSize uint32, pageToken string) (jobs []store.Job, nextPageToken string, err error) {
	return jobs, nextPageToken, nil
}

func (m *Memory) RetryDLQJob(ctx context.Context, jobID string) error {
	return nil
}

func (m *Memory) RetryDLQQueue(ctx context.Context, queueName string) (retriedCount uint32, err error) {
	return retriedCount, nil
}
