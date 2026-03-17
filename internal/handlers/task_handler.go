package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/workspace/repo/internal/models"
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GetAll handles GET /tasks.
func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks := h.store.GetAll()
	if tasks == nil {
		tasks = []models.Task{}
	}
	writeJSON(w, http.StatusOK, tasks)
}

// GetByID handles GET /tasks/{id}.
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, task)
}
