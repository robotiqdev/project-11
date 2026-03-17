package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/workspace/repo/internal/models"
)

// TaskStore defines the store interface used by TaskHandler.
type TaskStore interface {
	Create(req models.CreateTaskRequest) (models.Task, error)
	GetByID(id int64) (models.Task, error)
	Update(id int64, req models.UpdateTaskRequest) (models.Task, error)
}

// TaskHandler handles HTTP requests for tasks.
type TaskHandler struct {
	store TaskStore
}

// NewTaskHandler creates a new TaskHandler with the given store.
func NewTaskHandler(store TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Update handles PUT /tasks/{id} — not yet implemented.
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
}
