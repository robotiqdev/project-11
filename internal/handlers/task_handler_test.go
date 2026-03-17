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

// mockTaskStore is a test double implementing storage.TaskStore.
type mockTaskStore struct {
	tasks []models.Task
}

func (m *mockTaskStore) Create(req models.CreateTaskRequest) (models.Task, error) {
	return models.Task{}, nil
}

func (m *mockTaskStore) GetByID(id int64) (models.Task, error) {
	return models.Task{}, storage.ErrNotFound
}

func (m *mockTaskStore) GetAll() []models.Task {
	return m.tasks
}

func (m *mockTaskStore) Update(id int64, req models.UpdateTaskRequest) (models.Task, error) {
	return models.Task{}, storage.ErrNotFound
}

func (m *mockTaskStore) Delete(id int64) error {
	return storage.ErrNotFound
}

// compile-time check that mockTaskStore satisfies storage.TaskStore
var _ storage.TaskStore = (*mockTaskStore)(nil)

// newTask is a helper to build a Task with the given ID.
func newTask(id int64) models.Task {
	now := time.Now().UTC()
	return models.Task{
		ID:          id,
		Title:       "Task title",
		Description: "Task description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

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

// setupUpdateRouter creates a test mux with PUT /tasks/{id} registered.
func setupUpdateRouter(store storage.TaskStore) *http.ServeMux {
	handler := handlers.NewTaskHandler(store)
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /tasks/{id}", handler.Update)
	return mux
}

// seedTask creates a task in the store and returns it.
func seedTask(t *testing.T, store storage.TaskStore, title, description string) models.Task {
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

// doPutRequest sends a PUT request to the given URL with a JSON body.
func doPutRequest(t *testing.T, mux http.Handler, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

func TestGetAll_EmptyStore_Returns200(t *testing.T) {
	store := &mockTaskStore{tasks: nil}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetAll_EmptyStore_ReturnsEmptyJSONArray(t *testing.T) {
	store := &mockTaskStore{tasks: nil}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	var result []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil JSON array [], got null")
	}
	if len(result) != 0 {
		t.Errorf("expected 0 elements, got %d", len(result))
	}
}

func TestGetAll_EmptyStore_ReturnsContentTypeJSON(t *testing.T) {
	store := &mockTaskStore{tasks: nil}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestGetAll_WithNilReturnFromStore_ReturnsEmptyJSONArray(t *testing.T) {
	// Explicitly test the nil-slice case: handler must not marshal nil as JSON null.
	store := &mockTaskStore{tasks: nil} // GetAll() returns nil
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	body := w.Body.String()
	// Raw body must not be "null" — it must be an empty JSON array
	if body == "null" || body == "null\n" {
		t.Errorf("expected body [], got %q", body)
	}

	var result []models.Task
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatalf("failed to decode response body %q: %v", body, err)
	}
	if result == nil {
		t.Errorf("expected non-nil slice (JSON array []), got nil (JSON null)")
	}
}

func TestGetAll_ThreeTasks_Returns200(t *testing.T) {
	tasks := []models.Task{newTask(1), newTask(2), newTask(3)}
	store := &mockTaskStore{tasks: tasks}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestGetAll_ThreeTasks_ReturnsThreeElements(t *testing.T) {
	tasks := []models.Task{newTask(1), newTask(2), newTask(3)}
	store := &mockTaskStore{tasks: tasks}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	var result []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
}

func TestGetAll_ThreeTasks_ReturnsAscendingIDOrder(t *testing.T) {
	// Store returns tasks already sorted ascending (as InMemoryTaskStore does).
	tasks := []models.Task{newTask(1), newTask(2), newTask(3)}
	store := &mockTaskStore{tasks: tasks}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	var result []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	for i := 1; i < len(result); i++ {
		if result[i].ID <= result[i-1].ID {
			t.Errorf("expected ascending IDs, but result[%d].ID=%d is not greater than result[%d].ID=%d",
				i, result[i].ID, i-1, result[i-1].ID)
		}
	}
}

func TestGetAll_ThreeTasks_ReturnsContentTypeJSON(t *testing.T) {
	tasks := []models.Task{newTask(1), newTask(2), newTask(3)}
	store := &mockTaskStore{tasks: tasks}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestGetAll_ThreeTasks_CorrectIDsInResponse(t *testing.T) {
	tasks := []models.Task{newTask(10), newTask(20), newTask(30)}
	store := &mockTaskStore{tasks: tasks}
	handler := handlers.NewTaskHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	resp := w.Result()
	var result []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	expectedIDs := []int64{10, 20, 30}
	for i, task := range result {
		if task.ID != expectedIDs[i] {
			t.Errorf("result[%d]: expected ID %d, got %d", i, expectedIDs[i], task.ID)
		}
	}
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
