package storage

import (
	"errors"
	"sync"

	"github.com/workspace/repo/internal/models"
)

// ErrNotFound is returned when a task with the given ID does not exist.
var ErrNotFound = errors.New("task not found")

// IDGenerator is a stub for ID generation. Full implementation in a separate task.
type IDGenerator struct{}

// NewIDGenerator creates a new IDGenerator.
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// NextID returns the next available task ID. Full implementation in a separate task.
func (g *IDGenerator) NextID() int64 { return 0 }

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

// Create adds a new task to the store using the provided request. Stub implementation.
func (s *InMemoryTaskStore) Create(req models.CreateTaskRequest) (models.Task, error) {
	return models.Task{}, nil
}

// GetByID returns the task with the given ID, or ErrNotFound if it does not exist. Stub implementation.
func (s *InMemoryTaskStore) GetByID(id int64) (models.Task, error) {
	return models.Task{}, ErrNotFound
}

// GetAll returns all tasks currently in the store. Stub implementation.
func (s *InMemoryTaskStore) GetAll() []models.Task {
	return nil
}

// Update replaces the fields of the task identified by id using the provided request. Stub implementation.
func (s *InMemoryTaskStore) Update(id int64, req models.UpdateTaskRequest) (models.Task, error) {
	return models.Task{}, ErrNotFound
}

// Delete removes the task with the given ID from the store. Stub implementation.
func (s *InMemoryTaskStore) Delete(id int64) error {
	return ErrNotFound
}
