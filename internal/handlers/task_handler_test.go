package handlers_test

import (
	"encoding/json"
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
