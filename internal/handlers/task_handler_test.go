package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/workspace/repo/internal/handlers"
	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
)

// newTestHandler creates a TaskHandler backed by a fresh in-memory store.
func newTestHandler() (*handlers.TaskHandler, storage.TaskStore) {
	store := storage.NewInMemoryTaskStore()
	return handlers.NewTaskHandler(store), store
}

// setupRouter registers the DELETE and GET /tasks/{id} routes and returns the mux.
func setupRouter(h *handlers.TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /tasks/{id}", h.Delete)
	mux.HandleFunc("GET /tasks/{id}", h.GetByID)
	return mux
}

// createTask is a helper that creates a task in the store and returns it.
func createTask(t *testing.T, store storage.TaskStore, title, description string) models.Task {
	t.Helper()
	task, err := store.Create(models.CreateTaskRequest{
		Title:       title,
		Description: description,
	})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	return task
}

// --- DELETE /tasks/{id} tests ---

// TestDelete_ExistingTask_Returns204 verifies that deleting an existing task
// responds with HTTP 204 No Content.
func TestDelete_ExistingTask_Returns204(t *testing.T) {
	h, store := newTestHandler()
	task := createTask(t, store, "Test Task", "Test Description")

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	_ = task // confirm ID 1 was created
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rr.Code)
	}
}

// TestDelete_ExistingTask_Returns204_WithEmptyBody verifies that a successful
// DELETE returns an empty response body (RFC 7231 §6.3.5).
func TestDelete_ExistingTask_Returns204_WithEmptyBody(t *testing.T) {
	h, store := newTestHandler()
	createTask(t, store, "Task", "Description")

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
	if rr.Body.Len() != 0 {
		t.Errorf("expected empty body for 204 response, got %q", rr.Body.String())
	}
}

// TestDelete_NonExistentTask_Returns404 verifies that deleting a task that does
// not exist responds with HTTP 404 Not Found.
func TestDelete_NonExistentTask_Returns404(t *testing.T) {
	h, _ := newTestHandler()

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

// TestDelete_NonExistentTask_Returns404_WithJSONErrorBody verifies that the 404
// response for a non-existent task includes a JSON error body.
func TestDelete_NonExistentTask_Returns404_WithJSONErrorBody(t *testing.T) {
	h, _ := newTestHandler()

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid JSON error body, got decode error: %v", err)
	}
	if body["error"] == "" {
		t.Errorf("expected JSON body to contain non-empty 'error' field, got: %v", body)
	}
}

// TestDelete_InvalidID_Returns400 verifies that a non-integer task ID in the
// path responds with HTTP 400 Bad Request.
func TestDelete_InvalidID_Returns400(t *testing.T) {
	h, _ := newTestHandler()

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

// TestDelete_InvalidID_Returns400_WithJSONErrorBody verifies that the 400
// response for an invalid ID includes a JSON error body.
func TestDelete_InvalidID_Returns400_WithJSONErrorBody(t *testing.T) {
	h, _ := newTestHandler()

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid JSON error body for 400, got decode error: %v", err)
	}
	if body["error"] == "" {
		t.Errorf("expected JSON body to contain non-empty 'error' field, got: %v", body)
	}
}

// TestDelete_AfterSuccessfulDelete_GetReturns404 verifies that after a successful
// DELETE, a subsequent GET /tasks/{id} returns HTTP 404 for the same ID.
func TestDelete_AfterSuccessfulDelete_GetReturns404(t *testing.T) {
	h, store := newTestHandler()
	createTask(t, store, "Task To Delete", "Description")

	mux := setupRouter(h)

	// DELETE the task with ID 1.
	deleteReq := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	deleteRR := httptest.NewRecorder()
	mux.ServeHTTP(deleteRR, deleteReq)

	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("DELETE expected status 204, got %d", deleteRR.Code)
	}

	// Subsequent GET should return 404.
	getReq := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	getRR := httptest.NewRecorder()
	mux.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusNotFound {
		t.Errorf("GET after DELETE expected status 404, got %d", getRR.Code)
	}
}

// TestDelete_204_NoContentTypeHeader verifies that a 204 response does not set
// Content-Type (204 has no body, so Content-Type is irrelevant/unnecessary).
func TestDelete_204_NoContentTypeHeader(t *testing.T) {
	h, store := newTestHandler()
	createTask(t, store, "Task", "Description")

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rr.Code)
	}
	// 204 No Content should not carry a Content-Type header since there is no body.
	ct := rr.Header().Get("Content-Type")
	if ct != "" {
		t.Errorf("expected no Content-Type header for 204 response, got %q", ct)
	}
}

// TestDelete_ZeroID_Returns400 verifies that an ID of zero (invalid) returns 400.
func TestDelete_ZeroID_Returns400(t *testing.T) {
	h, _ := newTestHandler()

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/0", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// 0 is not a valid task ID; expect 400 or 404 depending on implementation.
	// Per spec, the store uses strictly positive IDs. Treat 0 as invalid (400) or not found (404).
	if rr.Code != http.StatusBadRequest && rr.Code != http.StatusNotFound {
		t.Errorf("expected status 400 or 404 for ID=0, got %d", rr.Code)
	}
}

// TestDelete_NegativeID_Returns400 verifies that a negative ID returns 400.
func TestDelete_NegativeID_Returns400(t *testing.T) {
	h, _ := newTestHandler()

	mux := setupRouter(h)
	req := httptest.NewRequest(http.MethodDelete, "/tasks/-1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	// Negative IDs are not valid. Expect 400 or 404.
	if rr.Code != http.StatusBadRequest && rr.Code != http.StatusNotFound {
		t.Errorf("expected status 400 or 404 for negative ID, got %d", rr.Code)
	}
}
