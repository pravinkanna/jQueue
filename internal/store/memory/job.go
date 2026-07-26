package memory

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/pravinkanna/jQueue/internal/store"
)

// Job
func (m *Memory) EnqueueJob(ctx context.Context, enqueueParams store.EnqueueParams) (jobID string, isDuplicate bool, err error) {
	qName := enqueueParams.Queue
	// Check whether the queue exist
	if _, ok := m.queueWithJobIDs[qName]; !ok {
		return "", false, store.ErrQueueNotFound
	}

	// Generate a new UUID v7
	id, err := uuid.NewV7()
	if err != nil {
		return "", false, err
	}
	jobID = id.String()

	// Create the Job
	m.jobs[jobID] = &store.Job{
		JobID:          jobID,
		IdempotencyKey: enqueueParams.IdempotencyKey,
		Queue:          enqueueParams.Queue,
		Payload:        enqueueParams.Payload,
		MaxRetries:     enqueueParams.MaxRetries,
		State:          store.JobStatePending,
		RetryCount:     0,
		CreatedAt:      time.Now(),
		ScheduledAt:    enqueueParams.RunAt,
	}
	m.queueWithJobIDs[qName].jobIDs = append(m.queueWithJobIDs[qName].jobIDs, jobID)
	m.queueWithJobIDs[qName].queue.PendingCount++

	return jobID, isDuplicate, nil
}

func (m *Memory) GetJob(ctx context.Context, jobID string) (job store.Job, err error) {
	// Check job exist and return the job
	j, ok := m.jobs[jobID]
	if !ok {
		return store.Job{}, store.ErrJobNotFound
	}
	job = *j
	return job, nil
}

func (m *Memory) CancelJob(ctx context.Context, jobID string) (state store.JobState, err error) {
	// Check job exist
	job, ok := m.jobs[jobID]
	if !ok {
		return store.JobStateUnspecified, store.ErrJobNotFound
	}

	// If the job is already started or done return error
	if job.State != store.JobStateScheduled && job.State != store.JobStatePending {
		return store.JobStateUnspecified, store.ErrJobAlreadyStatedOrCompleted
	}

	if job.State == store.JobStateScheduled {
		m.queueWithJobIDs[job.Queue].queue.ScheduledCount--
	} else {
		m.queueWithJobIDs[job.Queue].queue.PendingCount--
	}

	// Set job state to cancelled
	job.State = store.JobStateCancelled

	// Remove the job from jobIDs
	m.queueWithJobIDs[job.Queue].jobIDs = slices.DeleteFunc(
		m.queueWithJobIDs[job.Queue].jobIDs, func(x string) bool {
			return x == jobID
		})

	return job.State, nil
}

// TODO: Will do later
func (m *Memory) ListJobs(ctx context.Context, queue string, state store.JobState, pageSize uint32, pageToken string) (jobs []store.Job, nextPageToken string, err error) {
	return jobs, nextPageToken, nil
}

// TODO: Will do later after writing DLQ logic
func (m *Memory) RetryDLQJob(ctx context.Context, jobID string) error {
	return nil
}

// TODO: Will do later after writing DLQ logic
func (m *Memory) RetryDLQQueue(ctx context.Context, queueName string) (retriedCount uint32, err error) {
	return retriedCount, nil
}
