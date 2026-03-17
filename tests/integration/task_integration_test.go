package integration_test

import (
	"testing"

	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
)

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
