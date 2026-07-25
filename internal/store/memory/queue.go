package memory

import (
	"context"

	"github.com/pravinkanna/jQueue/internal/store"
)

// Queue
func (m *Memory) CreateQueue(ctx context.Context, name string) error {
	return nil
}

func (m *Memory) DeleteQueue(ctx context.Context, name string) error {
	return nil
}

func (m *Memory) PurgeQueue(ctx context.Context, name string) (purgedCount uint64, err error) {
	return purgedCount, nil
}

func (m *Memory) ListQueues(ctx context.Context) (queues []store.Queue, err error) {
	return queues, nil
}

func (m *Memory) GetQueueStatus(ctx context.Context, name string) (queue store.Queue, err error) {
	return queue, nil
}
