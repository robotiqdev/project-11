package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/workspace/repo/internal/handlers"
	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
)

// setupUpdateRouter creates a test mux with PUT /tasks/{id} registered.
func setupUpdateRouter(store *storage.TaskStore) *http.ServeMux {
	handler := handlers.NewTaskHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", handler.Update)
	return mux
}

// seedTask creates a task in the store and returns it.
func seedTask(t *testing.T, store *storage.TaskStore, title, description string) models.Task {
	t.Helper()
	task, err := store.Create(models.CreateTaskRequest{
		Title:       title,
		Description: description,
	})
	if err != nil {
		t.Fatalf("seedTask: Create returned unexpected error: %v", err)
	}
	return task
}

// doPutRequest sends a PUT request to the given URL with a JSON body.
func doPutRequest(t *testing.T, mux http.Handler, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

// --- Valid PUT /tasks/{id} tests ---

func TestUpdate_ValidRequest_Returns200(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	time.Sleep(1 * time.Millisecond)

	body := `{"title":"Updated Title","description":"Updated Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_ValidRequest_ResponseBodyHasUpdatedTitle(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	time.Sleep(1 * time.Millisecond)

	body := `{"title":"New Title","description":"New Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if resp.Title != "New Title" {
		t.Errorf("expected title %q in response, got %q", "New Title", resp.Title)
	}
}

func TestUpdate_ValidRequest_ResponseBodyHasUpdatedDescription(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	time.Sleep(1 * time.Millisecond)

	body := `{"title":"New Title","description":"New Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if resp.Description != "New Description" {
		t.Errorf("expected description %q in response, got %q", "New Description", resp.Description)
	}
}

func TestUpdate_ValidRequest_UpdatedAtIsAfterCreatedAt(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	time.Sleep(2 * time.Millisecond)

	body := `{"title":"New Title","description":"New Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if resp.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be set in response, got zero time")
	}
	if !resp.UpdatedAt.After(created.CreatedAt) {
		t.Errorf("expected updated_at (%v) to be after created_at (%v)", resp.UpdatedAt, created.CreatedAt)
	}
}

func TestUpdate_ValidRequest_ResponseBodyHasCorrectID(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"title":"New Title","description":"New Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if resp.ID != created.ID {
		t.Errorf("expected id %d in response, got %d", created.ID, resp.ID)
	}
}

func TestUpdate_ValidRequest_ResponseBodyPreservesCreatedAt(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"title":"New Title","description":"New Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if !resp.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("expected created_at to be preserved: got %v, want %v", resp.CreatedAt, created.CreatedAt)
	}
}

func TestUpdate_ValidRequest_ContentTypeIsJSON(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"title":"New Title","description":"New Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	ct := w.Header().Get("Content-Type")
	if ct == "" {
		t.Error("expected Content-Type header to be set")
	}
}

// --- PUT /tasks/999 non-existent task tests ---

func TestUpdate_NonExistentID_Returns404(t *testing.T) {
	store := storage.NewTaskStore()
	mux := setupUpdateRouter(store)

	body := `{"title":"Updated Title","description":"Updated Description"}`
	w := doPutRequest(t, mux, "/tasks/999", body)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404 for non-existent task, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_NonExistentID_ResponseBodyContainsError(t *testing.T) {
	store := storage.NewTaskStore()
	mux := setupUpdateRouter(store)

	body := `{"title":"Updated Title","description":"Updated Description"}`
	w := doPutRequest(t, mux, "/tasks/999", body)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response body: %v", err)
	}
	if _, ok := resp["error"]; !ok {
		t.Errorf("expected 'error' key in response body, got: %v", resp)
	}
}

// --- PUT /tasks/abc invalid id tests ---

func TestUpdate_InvalidIDFormat_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	mux := setupUpdateRouter(store)

	body := `{"title":"Updated Title","description":"Updated Description"}`
	w := doPutRequest(t, mux, "/tasks/abc", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_InvalidIDFormat_ResponseBodyContainsError(t *testing.T) {
	store := storage.NewTaskStore()
	mux := setupUpdateRouter(store)

	body := `{"title":"Updated Title","description":"Updated Description"}`
	w := doPutRequest(t, mux, "/tasks/abc", body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response body: %v", err)
	}
	if _, ok := resp["error"]; !ok {
		t.Errorf("expected 'error' key in response body, got: %v", resp)
	}
}

func TestUpdate_InvalidIDFormat_NonNumericString_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	mux := setupUpdateRouter(store)

	body := `{"title":"Updated Title","description":"Updated Description"}`
	w := doPutRequest(t, mux, "/tasks/not-a-number", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for non-numeric id, got %d; body: %s", w.Code, w.Body.String())
	}
}

// --- PUT with invalid JSON body tests ---

func TestUpdate_InvalidJSON_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"title": invalid json`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_InvalidJSON_ResponseBodyContainsError(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `not json at all`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response body: %v", err)
	}
	if _, ok := resp["error"]; !ok {
		t.Errorf("expected 'error' key in response body, got: %v", resp)
	}
}

func TestUpdate_EmptyBody_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), "")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty body, got %d; body: %s", w.Code, w.Body.String())
	}
}

// --- PUT with validation failure tests ---

func TestUpdate_EmptyTitle_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"title":"","description":"Some Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty title, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_EmptyTitle_ResponseBodyContainsError(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"title":"","description":"Some Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response body: %v", err)
	}
	if _, ok := resp["error"]; !ok {
		t.Errorf("expected 'error' key in response body, got: %v", resp)
	}
}

func TestUpdate_MissingTitle_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"description":"Some Description"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing title, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_EmptyDescription_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	body := `{"title":"Some Title","description":""}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty description, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestUpdate_TitleTooLong_Returns400(t *testing.T) {
	store := storage.NewTaskStore()
	created := seedTask(t, store, "Original Title", "Original Description")
	mux := setupUpdateRouter(store)

	longTitle := make([]byte, 201)
	for i := range longTitle {
		longTitle[i] = 'a'
	}
	body := fmt.Sprintf(`{"title":%q,"description":"Some Description"}`, string(longTitle))
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", created.ID), body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for title too long, got %d; body: %s", w.Code, w.Body.String())
	}
}

// --- Error ordering / priority tests ---

// Invalid ID format should be caught before JSON parsing.
func TestUpdate_InvalidIDFormat_DoesNotReturnTaskNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	mux := setupUpdateRouter(store)

	body := `{"title":"Updated Title","description":"Updated Description"}`
	w := doPutRequest(t, mux, "/tasks/xyz", body)

	if w.Code == http.StatusNotFound {
		t.Errorf("invalid id should return 400, not 404; got %d", w.Code)
	}
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id format, got %d", w.Code)
	}
}

// Valid update should not affect other tasks.
func TestUpdate_ValidRequest_DoesNotAffectOtherTasks(t *testing.T) {
	store := storage.NewTaskStore()
	task1 := seedTask(t, store, "Task One", "Description One")
	task2 := seedTask(t, store, "Task Two", "Description Two")
	mux := setupUpdateRouter(store)

	body := `{"title":"Task One Updated","description":"Description One Updated"}`
	w := doPutRequest(t, mux, fmt.Sprintf("/tasks/%d", task1.ID), body)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for valid update, got %d", w.Code)
	}

	// Verify task2 is unchanged in the store.
	got2, err := store.GetByID(task2.ID)
	if err != nil {
		t.Fatalf("GetByID task2 returned error: %v", err)
	}
	if got2.Title != task2.Title {
		t.Errorf("task2 title was mutated: got %q, want %q", got2.Title, task2.Title)
	}
}
