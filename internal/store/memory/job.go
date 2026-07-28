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

	// idempotency deduplication (if key given) - check a job there with the same idempotency key
	// TODO: Implement it as hashmap for o(1) check
	if enqueueParams.IdempotencyKey != "" {
		for _, j := range m.jobs {
			if enqueueParams.IdempotencyKey == j.IdempotencyKey && enqueueParams.Queue == j.Queue {
				return j.JobID, true, nil
			}
		}
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
		State:          store.JobStateScheduled,
		RetryCount:     0,
		CreatedAt:      time.Now(),
		ScheduledAt:    enqueueParams.RunAt,
	}
	m.queueWithJobIDs[qName].scheduledJobIDs = append(m.queueWithJobIDs[qName].scheduledJobIDs, jobID)

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
		return store.JobStateUnspecified, store.ErrJobAlreadyStartedOrCompleted
	}

	if job.State == store.JobStateScheduled {
		// Remove the job from jobIDs
		m.queueWithJobIDs[job.Queue].scheduledJobIDs = slices.DeleteFunc(
			m.queueWithJobIDs[job.Queue].scheduledJobIDs, func(x string) bool {
				return x == jobID
			})
	} else {
		// Remove the job from jobIDs
		m.queueWithJobIDs[job.Queue].pendingJobIDs = slices.DeleteFunc(
			m.queueWithJobIDs[job.Queue].pendingJobIDs, func(x string) bool {
				return x == jobID
			})
	}

	// Set job state to cancelled
	job.State = store.JobStateCancelled

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
