package memory

import (
	"sync"
	"time"

	"github.com/pravinkanna/jQueue/internal/store"
)

type Memory struct {
	mu              sync.Mutex
	queueWithJobIDs map[string]*queueWithJobIDs // k - queueName, V - queueMeta and jobIds
	jobs            map[string]*store.Job       // K - JobID, V - Job Data
	leases          map[string]*Lease           // K - LeaseToken, V - Lease Data
}

type queueWithJobIDs struct {
	queue           store.Queue
	dlqJobIDs       []string
	pendingJobIDs   []string
	scheduledJobIDs []string
}

type Lease struct {
	jobID     string
	expiresAt time.Time
}

func New() *Memory {
	return &Memory{
		queueWithJobIDs: make(map[string]*queueWithJobIDs),
		jobs:            make(map[string]*store.Job),
		leases:          make(map[string]*Lease),
	}
}
