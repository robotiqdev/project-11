package storage

import (
	"errors"
	"sort"
	"sync"
	"time"

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
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.idGen.NextID()
	now := time.Now().UTC()
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

// GetByID returns the task with the given ID or ErrNotFound.
func (s *InMemoryTaskStore) GetByID(id int64) (models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	return task, nil
}

// GetAll returns all stored tasks sorted by ID.
func (s *InMemoryTaskStore) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	return tasks
}

// Update modifies an existing task or returns ErrNotFound.
func (s *InMemoryTaskStore) Update(id int64, req models.UpdateTaskRequest) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, ErrNotFound
	}
	task.Title = req.Title
	task.Description = req.Description
	task.UpdatedAt = time.Now().UTC()
	s.tasks[id] = task
	return task, nil
}

// Delete removes a task by ID or returns ErrNotFound.
// IDs are never reused: the id generator is monotonically increasing and delete only removes the map entry.
func (s *InMemoryTaskStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
