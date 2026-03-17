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
	s.mu.Lock()
	defer s.mu.Unlock()
	task := models.Task{
		ID:          s.nextID,
		Title:       req.Title,
		Description: req.Description,
	}
	s.tasks[s.nextID] = task
	s.nextID++
	return task, nil
}

// GetByID retrieves a task by its ID.
func (s *TaskStore) GetByID(id int64) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	return task, nil
}

// Update replaces a task's fields by ID.
func (s *TaskStore) Update(id int64, req models.UpdateTaskRequest) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	task.Title = req.Title
	task.Description = req.Description
	s.tasks[id] = task
	return task, nil
}

// Delete removes a task by ID.
func (s *TaskStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.tasks[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
