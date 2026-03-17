package handlers

import (
	"net/http"

	"github.com/workspace/repo/internal/storage"
)

// TaskHandler handles HTTP requests for task operations.
type TaskHandler struct {
	store storage.TaskStore
}

// NewTaskHandler creates a new TaskHandler with the given store.
func NewTaskHandler(store storage.TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

// Delete handles DELETE /tasks/{id} requests.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// stub — implementation pending
}

// GetByID handles GET /tasks/{id} requests.
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// stub — implementation pending
}
