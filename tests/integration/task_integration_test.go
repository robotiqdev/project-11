package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/workspace/repo/internal/handlers"
	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/router"
	"github.com/workspace/repo/internal/storage"
)

// newIntegrationRouter creates a router backed by a fresh in-memory store for
// integration tests. It returns both the handler (for direct store access) and
// the routed http.Handler.
func newIntegrationRouter() (http.Handler, *handlers.TaskHandler, *storage.InMemoryTaskStore) {
	store := storage.NewInMemoryTaskStore()
	h := handlers.NewTaskHandler(store)
	r := router.New(h)
	return r, h, store
}

// newIntegrationStore returns a fresh InMemoryTaskStore for integration tests.
func newIntegrationStore() *storage.InMemoryTaskStore {
	return storage.NewInMemoryTaskStore()
}

// TestIDCounter_NewIDGreaterThanDeletedID verifies that after creating and
// deleting a task, the next created task receives an ID strictly greater than
// the deleted task's ID — confirming IDs are never reused.
func TestIDCounter_NewIDGreaterThanDeletedID(t *testing.T) {
	s := newIntegrationStore()

	// POST (create) a task.
	created, err := s.Create(models.CreateTaskRequest{
		Title:       "Integration Task",
		Description: "Integration Description",
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	deletedID := created.ID

	// DELETE the task.
	if err := s.Delete(deletedID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// POST (create) a new task.
	next, err := s.Create(models.CreateTaskRequest{
		Title:       "Next Task",
		Description: "Next Description",
	})
	if err != nil {
		t.Fatalf("create next failed: %v", err)
	}

	// Assert new ID is strictly greater than the deleted ID.
	if next.ID <= deletedID {
		t.Errorf("expected next task ID (%d) > deleted task ID (%d): IDs must not be reused", next.ID, deletedID)
	}
}

// TestIDCounter_NeverReusesAnyDeletedID verifies that after creating multiple
// tasks and deleting some, subsequent tasks always receive new monotonically
// increasing IDs rather than recycling deleted ones.
func TestIDCounter_NeverReusesAnyDeletedID(t *testing.T) {
	s := newIntegrationStore()

	// Create three tasks.
	var ids [3]int64
	for i := 0; i < 3; i++ {
		task, err := s.Create(models.CreateTaskRequest{
			Title:       "Task",
			Description: "Description",
		})
		if err != nil {
			t.Fatalf("create task %d failed: %v", i, err)
		}
		ids[i] = task.ID
	}

	// Delete the first and second tasks.
	for _, id := range ids[:2] {
		if err := s.Delete(id); err != nil {
			t.Fatalf("delete task (id=%d) failed: %v", id, err)
		}
	}

	// Create two more tasks.
	deletedSet := map[int64]bool{ids[0]: true, ids[1]: true}
	maxDeleted := ids[1]
	if ids[0] > ids[1] {
		maxDeleted = ids[0]
	}

	for i := 0; i < 2; i++ {
		task, err := s.Create(models.CreateTaskRequest{
			Title:       "New Task",
			Description: "New Description",
		})
		if err != nil {
			t.Fatalf("create new task %d failed: %v", i, err)
		}
		if deletedSet[task.ID] {
			t.Errorf("new task %d received a reused (deleted) ID %d", i, task.ID)
		}
		if task.ID <= maxDeleted {
			t.Errorf("new task %d ID (%d) is not greater than max deleted ID (%d)", i, task.ID, maxDeleted)
		}
	}
}

// --- Router route registration tests ---

// TestRouter_DeleteTaskByID_IsRouted verifies that the router registers the
// DELETE /tasks/{id} route, so a DELETE request is dispatched to the handler
// rather than returning 404 from the router itself. A pre-seeded task is used
// so the handler can respond with 204 No Content (confirming the route exists).
func TestRouter_DeleteTaskByID_IsRouted(t *testing.T) {
	r, _, store := newIntegrationRouter()

	// Seed a task so the handler returns 204, not 404.
	task, err := store.Create(models.CreateTaskRequest{
		Title:       "Router Test Task",
		Description: "Router Test Description",
	})
	if err != nil {
		t.Fatalf("failed to seed task: %v", err)
	}

	path := "/tasks/" + itoa(task.ID)
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// 404 means the router has no matching route — the route is not registered.
	if rr.Code == http.StatusNotFound {
		t.Errorf("DELETE %s returned 404: DELETE /tasks/{id} route is not registered in the router", path)
	}
	// Expect 204 No Content when the task exists and the route is properly registered.
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", rr.Code)
	}
}

// TestRouter_PostTaskByID_Returns405 verifies that POST /tasks/{id} returns
// HTTP 405 Method Not Allowed. Go 1.22 method-prefixed patterns automatically
// respond 405 when a path matches but the HTTP method does not.
func TestRouter_PostTaskByID_Returns405(t *testing.T) {
	r, _, _ := newIntegrationRouter()

	req := httptest.NewRequest(http.MethodPost, "/tasks/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405 Method Not Allowed for POST /tasks/1, got %d", rr.Code)
	}
}

// itoa converts an int64 to its decimal string representation.
func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
