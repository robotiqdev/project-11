package storage

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/workspace/repo/internal/models"
)

// ErrNotFound is returned when a task with the given ID does not exist.
var ErrNotFound = errors.New("task not found")

// IDGenerator generates unique sequential IDs.
type IDGenerator struct {
	counter int64
}

// NewIDGenerator creates a new IDGenerator.
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

// NextID returns the next unique task ID.
func (g *IDGenerator) NextID() int64 {
	return atomic.AddInt64(&g.counter, 1)
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

// Create adds a new task to the store using the provided request.
func (s *InMemoryTaskStore) Create(req models.CreateTaskRequest) (models.Task, error) {
	if err := req.Validate(); err != nil {
		return models.Task{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.idGen.NextID()
	now := time.Now()
	task := models.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[id] = task
	return task, nil
}

// GetByID returns the task with the given ID, or ErrNotFound if it does not exist.
func (s *InMemoryTaskStore) GetByID(id int64) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	return task, nil
}

// GetAll returns all tasks currently in the store.
func (s *InMemoryTaskStore) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

// Update replaces the fields of the task identified by id using the provided request.
func (s *InMemoryTaskStore) Update(id int64, req models.UpdateTaskRequest) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}

	task.Title = req.Title
	task.Description = req.Description
	task.UpdatedAt = time.Now()
	s.tasks[id] = task
	return task, nil
}

// Delete removes the task with the given ID from the store.
func (s *InMemoryTaskStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
