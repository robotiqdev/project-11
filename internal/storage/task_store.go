package storage

import (
	"errors"
	"sync"

	"github.com/workspace/repo/internal/models"
)

// ErrNotFound is returned when a task is not found in the store.
var ErrNotFound = errors.New("task not found")

// TaskStore is an in-memory store for tasks.
type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[int64]models.Task
	nextID int64
}

// NewTaskStore creates a new TaskStore.
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int64]models.Task),
		nextID: 1,
	}
}

// Create stores a new task and returns it with an assigned ID.
func (s *TaskStore) Create(req models.CreateTaskRequest) (models.Task, error) {
	return models.Task{}, nil
}

// GetByID retrieves a task by its ID.
func (s *TaskStore) GetByID(id int64) (models.Task, error) {
	return models.Task{}, nil
}

// Update replaces a task's fields by ID.
func (s *TaskStore) Update(id int64, req models.UpdateTaskRequest) (models.Task, error) {
	return models.Task{}, nil
}

// Delete removes a task by ID.
func (s *TaskStore) Delete(id int64) error {
	return nil
}
