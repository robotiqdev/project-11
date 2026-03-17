package handlers

import (
	"net/http"

	"github.com/workspace/repo/internal/storage"
)

// TaskHandler handles HTTP requests for task resources.
type TaskHandler struct {
	store storage.TaskStore
}

// NewTaskHandler creates a new TaskHandler with the given store.
func NewTaskHandler(store storage.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

// GetAll handles GET /tasks — stub, not yet implemented.
func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// not implemented
}
