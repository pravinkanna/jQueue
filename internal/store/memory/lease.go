package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Lease
func (m *Memory) LeaseJob(ctx context.Context, queue string, leaseDuration time.Duration) (leasedJob store.LeasedJob, err error) {
	// check the queue exist
	if _, ok := m.queueWithJobIDs[queue]; !ok {
		return store.LeasedJob{}, store.ErrQueueNotFound
	}

	// make sure the queue is not empty
	pendingIdsCount := len(m.queueWithJobIDs[queue].pendingJobIDs)
	if pendingIdsCount == 0 {
		// TODO: Implement Long polling
		return store.LeasedJob{}, store.ErrQueueEmpty
	}

	// FIFO - queue
	jobID := m.queueWithJobIDs[queue].pendingJobIDs[0]
	m.queueWithJobIDs[queue].pendingJobIDs = m.queueWithJobIDs[queue].pendingJobIDs[1:]

	// Generate lease token and store it in job struct
	id, err := uuid.NewV7()
	if err != nil {
		return store.LeasedJob{}, err
	}
	leaseToken := id.String()
	expiresAt := time.Now().Add(leaseDuration)

	// change the job status to processing and set lease token and expiry time
	m.jobs[jobID].State = store.JobStateLeased
	m.leases[leaseToken] = &Lease{jobID: jobID, expiresAt: expiresAt}

	leasedJob = store.LeasedJob{
		Job:        *m.jobs[jobID],
		LeaseToken: leaseToken,
		ExpiresAt:  expiresAt,
	}

	m.queueWithJobIDs[queue].queue.LeasedCount++

	// return leased job
	return leasedJob, nil
}

func (m *Memory) ExtendJobLease(ctx context.Context, leaseToken string, duration time.Duration) (leaseExpiresAt time.Time, err error) {
	// Make sure the lease is available
	lease, ok := m.leases[leaseToken]
	if !ok {
		return time.Time{}, store.ErrLeaseNotFound
	}

	// Make sure it is not expired
	if time.Now().After(lease.expiresAt) {
		return time.Time{}, store.ErrLeaseExpired
	}

	// Increase the expiry time
	lease.expiresAt = lease.expiresAt.Add(duration)

	return lease.expiresAt, nil
}

func (m *Memory) AckJob(ctx context.Context, leaseToken string) error {
	// check the lease exist
	lease, ok := m.leases[leaseToken]
	if !ok {
		return store.ErrLeaseNotFound
	}
	jobID := lease.jobID

	// Make sure it is not expired
	if time.Now().After(lease.expiresAt) {
		return store.ErrLeaseExpired
	}

	// Make sure the jobId exist
	job, ok := m.jobs[jobID]
	if !ok {
		return store.ErrJobNotFound
	}

	// Make sure the job status is JobStateLeased
	if job.State != store.JobStateLeased {
		return store.ErrJobNotLeased
	}

	// Mark it as JobStateCompleted
	job.State = store.JobStateCompleted
	job.CompletedAt = time.Now()
	m.queueWithJobIDs[job.Queue].queue.CompletedCount++
	m.queueWithJobIDs[job.Queue].queue.LeasedCount--

	// remove it from the queue slice
	delete(m.leases, leaseToken)

	return nil
}

func (m *Memory) NackJob(ctx context.Context, leaseToken string, reason string) error {
	// check the lease exist
	lease, ok := m.leases[leaseToken]
	if !ok {
		return store.ErrLeaseNotFound
	}
	jobID := lease.jobID

	// Make sure it is not expired
	if time.Now().After(lease.expiresAt) {
		return store.ErrLeaseExpired
	}

	// Make sure the jobId exist
	job, ok := m.jobs[jobID]
	if !ok {
		return store.ErrJobNotFound
	}

	// Make sure the job status is JobStateLeased
	if job.State != store.JobStateLeased {
		return store.ErrJobNotLeased
	}

	job.LastError = reason
	job.RetryCount++

	// If the retry count > MaxRetries send to DLQ
	if job.RetryCount > job.MaxRetries {
		m.queueWithJobIDs[job.Queue].dlqJobIDs = append(m.queueWithJobIDs[job.Queue].dlqJobIDs, jobID)
		job.State = store.JobStateDLQ
		job.CompletedAt = time.Now()
	} else {
		// create delay with with exponential backoff
		baseDelay := 10 * time.Second
		delay := baseDelay << job.RetryCount // TODO: Add jitter later
		job.ScheduledAt = time.Now().Add(delay)
		m.queueWithJobIDs[job.Queue].scheduledJobIDs = append(m.queueWithJobIDs[job.Queue].scheduledJobIDs, jobID)
		job.State = store.JobStateScheduled
	}

	// remove it from the queue slice
	delete(m.leases, leaseToken)
	m.queueWithJobIDs[job.Queue].queue.LeasedCount--

	return nil
}
