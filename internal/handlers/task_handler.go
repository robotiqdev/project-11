package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"

	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
)

// TaskStore defines the data-access interface for task operations.
type TaskStore interface {
	Create(req models.CreateTaskRequest) (*models.Task, error)
	GetAll() ([]*models.Task, error)
	GetByID(id int64) (*models.Task, error)
	Update(id int64, req models.UpdateTaskRequest) (*models.Task, error)
	Delete(id int64) error
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

// GetAll handles GET /tasks requests.
func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// GetByID handles GET /tasks/{id} requests.
func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	task, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Update handles PUT /tasks/{id} requests.
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	task, err := h.store.Update(id, req)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Delete handles DELETE /tasks/{id} requests.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}
	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseID extracts the {id} path parameter and parses it as int64.
func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}
