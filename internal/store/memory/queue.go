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
	queue := store.Queue{
		Name: name,
	}
	jobIDs := []string{}
	m.queueWithJobIDs[name] = &queueWithJobIDs{
		queue:  queue,
		jobIDs: jobIDs,
	}

	return nil
}

func (m *Memory) DeleteQueue(ctx context.Context, name string) error {
	// Check whether the queue exist
	if _, ok := m.queueWithJobIDs[name]; !ok {
		return store.ErrQueueNotFound
	}

	// Make sure the queue is empty
	jobIDs := m.queueWithJobIDs[name].jobIDs
	if len(jobIDs) != 0 {
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
	jobIDs := m.queueWithJobIDs[name].jobIDs

	// Iterate through the jobIds and delete the job from job map
	for _, jobID := range jobIDs {
		delete(m.jobs, jobID)
	}

	// Get length of array
	purgedCount = uint64(len(jobIDs))

	// Make jobIDs slice empty
	m.queueWithJobIDs[name].jobIDs = []string{}

	return purgedCount, nil
}

func (m *Memory) ListQueues(ctx context.Context) (queues []store.Queue, err error) {
	queues = []store.Queue{}
	for _, qData := range m.queueWithJobIDs {
		queue := qData.queue
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
	return queue, nil
}
