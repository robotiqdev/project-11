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
