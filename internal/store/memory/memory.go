package memory

import (
	"sync"

	"github.com/pravinkanna/jQueue/internal/store"
)

type Memory struct {
	mu              sync.Mutex
	queueWithJobIDs map[string]*queueWithJobIDs // k - queueName, V - queueMeta and jobIds
	jobs            map[string]*store.Job       // K - JobID, V - Job Data
}

type queueWithJobIDs struct {
	queue  store.Queue
	jobIDs []string
}

func New() *Memory {
	return &Memory{
		queueWithJobIDs: make(map[string]*queueWithJobIDs),
		jobs:            make(map[string]*store.Job),
	}
}
