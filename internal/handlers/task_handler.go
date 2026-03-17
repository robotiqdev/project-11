package handlers

import (
	"encoding/json"
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
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	task, err := h.store.Create(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
