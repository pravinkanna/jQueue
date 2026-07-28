package memory

import (
	"context"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Queue
func (m *Memory) CreateQueue(ctx context.Context, name string) error {
	// Make sure the queue doesn't exist
	if _, ok := m.queueWithJobIDs[name]; ok {
		return store.ErrQueueExists
	}

	// Create the queue
	m.queueWithJobIDs[name] = &queueWithJobIDs{
		queue:           store.Queue{Name: name},
		dlqJobIDs:       []string{},
		pendingJobIDs:   []string{},
		scheduledJobIDs: []string{},
	}

	return nil
}

func (m *Memory) DeleteQueue(ctx context.Context, name string) error {
	// Check whether the queue exist
	if _, ok := m.queueWithJobIDs[name]; !ok {
		return store.ErrQueueNotFound
	}

	// Make sure the queue is empty
	pendingJobIDs := m.queueWithJobIDs[name].pendingJobIDs
	scheduledJobIDs := m.queueWithJobIDs[name].scheduledJobIDs
	dlqJobIDs := m.queueWithJobIDs[name].dlqJobIDs
	leasedCount := m.queueWithJobIDs[name].queue.LeasedCount
	if len(pendingJobIDs) != 0 || len(scheduledJobIDs) != 0 || len(dlqJobIDs) != 0 || leasedCount != 0 {
		return store.ErrQueueNotEmpty
	}

	// Delete the queue
	delete(m.queueWithJobIDs, name)

	return nil
}

func (m *Memory) PurgeQueue(ctx context.Context, name string) (purgedCount uint64, err error) {
	// Check whether the queue exist
	if _, ok := m.queueWithJobIDs[name]; !ok {
		return 0, store.ErrQueueNotFound
	}
	pendingJobIDs := m.queueWithJobIDs[name].pendingJobIDs
	scheduledJobIDs := m.queueWithJobIDs[name].scheduledJobIDs
	dlqJobIDs := m.queueWithJobIDs[name].dlqJobIDs

	// Iterate through the jobIds and delete the job from job map
	for _, jobID := range pendingJobIDs {
		delete(m.jobs, jobID)
	}

	// Iterate through the jobIds and delete the job from job map
	for _, jobID := range scheduledJobIDs {
		delete(m.jobs, jobID)
	}

	// Iterate through the jobIds and delete the job from job map
	for _, jobID := range dlqJobIDs {
		delete(m.jobs, jobID)
	}

	// Get length of array
	purgedCount = 0
	purgedCount += uint64(len(pendingJobIDs))
	purgedCount += uint64(len(scheduledJobIDs))
	purgedCount += uint64(len(dlqJobIDs))

	// Make jobIDs slice empty
	m.queueWithJobIDs[name].pendingJobIDs = []string{}
	m.queueWithJobIDs[name].scheduledJobIDs = []string{}
	m.queueWithJobIDs[name].dlqJobIDs = []string{}
	m.queueWithJobIDs[name].queue.CompletedCount = 0

	return purgedCount, nil
}

func (m *Memory) ListQueues(ctx context.Context) (queues []store.Queue, err error) {
	queues = []store.Queue{}
	for _, qData := range m.queueWithJobIDs {
		queue := qData.queue
		queue.ScheduledCount = uint64(len(qData.scheduledJobIDs))
		queue.PendingCount = uint64(len(qData.pendingJobIDs))
		queue.DLQCount = uint64(len(qData.dlqJobIDs))
		queues = append(queues, queue)
	}
	return queues, nil
}

func (m *Memory) GetQueueStatus(ctx context.Context, name string) (store.Queue, error) {
	queueWithJobIDs, ok := m.queueWithJobIDs[name]
	if !ok {
		return store.Queue{}, store.ErrQueueNotFound
	}

	queue := queueWithJobIDs.queue
	queue.ScheduledCount = uint64(len(queueWithJobIDs.scheduledJobIDs))
	queue.PendingCount = uint64(len(queueWithJobIDs.pendingJobIDs))
	queue.DLQCount = uint64(len(queueWithJobIDs.dlqJobIDs))
	return queue, nil
}
