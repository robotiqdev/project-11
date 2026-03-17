package storage

import (
	"errors"
	"sync"

	"github.com/workspace/repo/internal/models"
)

// ErrNotFound is returned when a task with the requested ID does not exist.
var ErrNotFound = errors.New("task not found")

// IDGenerator generates strictly increasing integer IDs.
type IDGenerator struct {
	mu      sync.Mutex
	current int64
}

// NextID returns the next unique ID.
func (g *IDGenerator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.current++
	return g.current
}

// TaskStore defines the interface for task persistence operations.
type TaskStore interface {
	Create(req models.CreateTaskRequest) (models.Task, error)
	GetByID(id int64) (models.Task, error)
	GetAll() []models.Task
	Update(id int64, req models.UpdateTaskRequest) (models.Task, error)
	Delete(id int64) error
}

// InMemoryTaskStore is a thread-safe in-memory implementation of TaskStore.
type InMemoryTaskStore struct {
	mu    sync.RWMutex
	tasks map[int64]models.Task
	idGen *IDGenerator
}

// NewInMemoryTaskStore creates and returns a new InMemoryTaskStore.
func NewInMemoryTaskStore() *InMemoryTaskStore {
	return &InMemoryTaskStore{
		tasks: make(map[int64]models.Task),
		idGen: &IDGenerator{},
	}
}

// Create stores a new task and returns it with a populated ID and timestamps.
func (s *InMemoryTaskStore) Create(req models.CreateTaskRequest) (models.Task, error) {
	return models.Task{}, errors.New("not implemented")
}

// GetByID returns the task with the given ID or ErrNotFound.
func (s *InMemoryTaskStore) GetByID(id int64) (models.Task, error) {
	return models.Task{}, errors.New("not implemented")
}

// GetAll returns all stored tasks.
func (s *InMemoryTaskStore) GetAll() []models.Task {
	return nil
}

// Update modifies an existing task or returns ErrNotFound.
func (s *InMemoryTaskStore) Update(id int64, req models.UpdateTaskRequest) (models.Task, error) {
	return models.Task{}, errors.New("not implemented")
}

// Delete removes a task by ID or returns ErrNotFound.
func (s *InMemoryTaskStore) Delete(id int64) error {
	return errors.New("not implemented")
}
