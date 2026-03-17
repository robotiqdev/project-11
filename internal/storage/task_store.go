package storage

import (
	"sync"

	"github.com/workspace/repo/internal/models"
)

// IDGenerator is a stub for ID generation. Full implementation in a separate task.
type IDGenerator struct{}

// NewIDGenerator creates a new IDGenerator.
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// InMemoryTaskStore holds tasks in memory.
type InMemoryTaskStore struct {
	mu    sync.RWMutex
	tasks map[int64]models.Task
	idGen *IDGenerator
}

// NewInMemoryTaskStore returns an InMemoryTaskStore with an initialized tasks map and IDGenerator.
func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{
		tasks: make(map[int64]models.Task),
		idGen: NewIDGenerator(),
	}
}
