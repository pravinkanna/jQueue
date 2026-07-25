package memory

import (
	"context"
	"fmt"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Queue
func (m *Memory) CreateQueue(ctx context.Context, name string) error {
	// Make sure the queue doesn't exist
	if _, ok := m.queues[name]; ok {
		return store.ErrQueueExists
	}

	// Create the queue
	m.queues[name] = store.Queue{
		Name: name,
	}

	// Create the Jobs slice for queue
	m.jobs[name] = []store.Job{}

	fmt.Println("Queue Created", m.queues[name])
	return nil
}

func (m *Memory) DeleteQueue(ctx context.Context, name string) error {
	// Check whether the queue exist
	if _, ok := m.queues[name]; !ok {
		return store.ErrQueueNotFound
	}

	// Make sure the queue is empty
	jobs := m.jobs[name]
	if len(jobs) != 0 {
		return store.ErrQueueNotEmpty
	}

	// Delete the queue
	delete(m.queues, name)
	delete(m.jobs, name)

	fmt.Println("Queue Deleted", m.queues[name])

	return nil
}

func (m *Memory) PurgeQueue(ctx context.Context, name string) (purgedCount uint64, err error) {
	// Check whether the queue exist
	if _, ok := m.jobs[name]; !ok {
		return 0, store.ErrQueueNotFound
	}

	// Get length of array
	purgedCount = uint64(len(m.jobs[name]))

	fmt.Println("purgedCount", len(m.jobs[name]))

	// Make array
	m.jobs[name] = []store.Job{}

	return purgedCount, nil
}

func (m *Memory) ListQueues(ctx context.Context) (queues []store.Queue, err error) {
	queues = []store.Queue{}
	for _, qData := range m.queues {
		queues = append(queues, qData)
	}
	return queues, nil
}

func (m *Memory) GetQueueStatus(ctx context.Context, name string) (store.Queue, error) {
	queue, ok := m.queues[name]
	if !ok {
		return store.Queue{}, store.ErrQueueNotFound
	}
	return queue, nil
}
