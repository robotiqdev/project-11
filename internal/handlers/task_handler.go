package handlers

import (
	"net/http"

	"github.com/workspace/repo/internal/models"
)

// TaskStore defines the data-access interface for task operations.
type TaskStore interface {
	Create(req models.CreateTaskRequest) (*models.Task, error)
}

// TaskHandler handles HTTP requests for task resources.
type TaskHandler struct {
	store TaskStore
}

// NewTaskHandler returns a new TaskHandler backed by the given store.
func NewTaskHandler(store TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

// Create handles POST /tasks requests.
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Not implemented yet.
}
