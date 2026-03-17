package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
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

// Update handles PUT /tasks/{id}.
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.store.Update(id, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}
