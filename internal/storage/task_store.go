package storage

import (
	"sync"

	"github.com/workspace/repo/internal/models"
)

// IDGenerator is a stub for ID generation. Full implementation in a separate task.
type IDGenerator struct{}

// NewIDGenerator creates a new IDGenerator stub.
func NewIDGenerator() *IDGenerator {
	return nil // stub: not implemented
}

// InMemoryTaskStore holds tasks in memory.
type InMemoryTaskStore struct {
	mu    sync.RWMutex
	tasks map[int64]models.Task
	idGen *IDGenerator
}

// NewInMemoryTaskStore returns a stub InMemoryTaskStore with uninitialized fields.
// Proper implementation must initialize tasks map and idGen.
func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{} // stub: tasks is nil, idGen is nil — tests should fail
}
