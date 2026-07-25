package memory

import (
	"sync"

	"github.com/pravinkanna/jQueue/internal/store"
)

type Memory struct {
	mu     sync.Mutex
	queues map[string]store.Queue // k - queueName, V - queueInfo
	jobs   map[string][]store.Job // K - queueName, V - list of jobs
}

func New() *Memory {
	return &Memory{
		queues: make(map[string]store.Queue),
		jobs:   make(map[string][]store.Job),
	}
}
