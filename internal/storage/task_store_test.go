package storage_test

import (
	"errors"
	"testing"

	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
)

// --- GetByID tests ---

func TestGetByID_NonExistentID_ReturnsError(t *testing.T) {
	store := storage.NewTaskStore()
	_, err := store.GetByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent ID, got nil")
	}
}

func TestGetByID_NonExistentID_ReturnsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	_, err := store.GetByID(999)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected errors.Is(err, storage.ErrNotFound) to be true, got err=%v", err)
	}
}

func TestGetByID_NonExistentID_ErrorsIsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	_, err := store.GetByID(42)
	if err == nil {
		t.Fatal("expected non-nil error for non-existent task")
	}
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected errors.Is(err, storage.ErrNotFound)==true, got err=%v", err)
	}
}

func TestGetByID_AfterCreate_ReturnsCorrectTask(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.CreateTaskRequest{
		Title:       "My Task",
		Description: "My Description",
	}
	created, err := store.Create(req)
	if err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("expected Create to assign a non-zero ID")
	}

	got, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID returned unexpected error after Create: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID mismatch: got %d, want %d", got.ID, created.ID)
	}
	if got.Title != req.Title {
		t.Errorf("Title mismatch: got %q, want %q", got.Title, req.Title)
	}
	if got.Description != req.Description {
		t.Errorf("Description mismatch: got %q, want %q", got.Description, req.Description)
	}
}

func TestGetByID_AfterCreate_ReturnsTaskWithMatchingTitle(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.CreateTaskRequest{
		Title:       "Lookup Title",
		Description: "Lookup Description",
	}
	created, err := store.Create(req)
	if err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}

	got, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID returned unexpected error: %v", err)
	}
	if got.Title != "Lookup Title" {
		t.Errorf("expected title %q, got %q", "Lookup Title", got.Title)
	}
}

func TestGetByID_AfterCreate_ReturnsTaskWithMatchingDescription(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.CreateTaskRequest{
		Title:       "Title",
		Description: "Expected Description",
	}
	created, err := store.Create(req)
	if err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}

	got, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID returned unexpected error: %v", err)
	}
	if got.Description != "Expected Description" {
		t.Errorf("expected description %q, got %q", "Expected Description", got.Description)
	}
}

func TestGetByID_IDZero_ReturnsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	_, err := store.GetByID(0)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for ID=0, got: %v", err)
	}
}

func TestGetByID_NegativeID_ReturnsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	_, err := store.GetByID(-1)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for negative ID, got: %v", err)
	}
}

// --- Update tests ---

func TestUpdate_NonExistentID_ReturnsError(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.UpdateTaskRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	_, err := store.Update(999, req)
	if err == nil {
		t.Fatal("expected error for Update with non-existent ID, got nil")
	}
}

func TestUpdate_NonExistentID_ReturnsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.UpdateTaskRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	_, err := store.Update(999, req)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected errors.Is(err, storage.ErrNotFound) to be true, got err=%v", err)
	}
}

func TestUpdate_NonExistentID_ErrorsIsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.UpdateTaskRequest{
		Title:       "Title",
		Description: "Description",
	}
	_, err := store.Update(42, req)
	if err == nil {
		t.Fatal("expected non-nil error for Update with non-existent task")
	}
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected errors.Is(err, storage.ErrNotFound)==true, got err=%v", err)
	}
}

func TestUpdate_IDZero_ReturnsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.UpdateTaskRequest{
		Title:       "Title",
		Description: "Description",
	}
	_, err := store.Update(0, req)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for Update with ID=0, got: %v", err)
	}
}

// --- Delete tests ---

func TestDelete_NonExistentID_ReturnsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	err := store.Delete(999)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("expected errors.Is(err, storage.ErrNotFound) to be true for Delete, got err=%v", err)
	}
}

func TestDelete_NonExistentID_ReturnsError(t *testing.T) {
	store := storage.NewTaskStore()
	err := store.Delete(999)
	if err == nil {
		t.Fatal("expected error for Delete with non-existent ID, got nil")
	}
}

// --- ErrNotFound sentinel tests ---

func TestErrNotFound_IsSentinelValue(t *testing.T) {
	if storage.ErrNotFound == nil {
		t.Fatal("storage.ErrNotFound must not be nil")
	}
	if storage.ErrNotFound.Error() == "" {
		t.Fatal("storage.ErrNotFound must have a non-empty error message")
	}
}

func TestErrNotFound_ErrorsIsComparison_GetByID(t *testing.T) {
	store := storage.NewTaskStore()
	_, err := store.GetByID(1234)
	if err == nil {
		t.Fatal("expected non-nil error for non-existent ID")
	}
	// errors.Is must work — this requires direct sentinel or wrapping via fmt.Errorf("%w")
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("errors.Is(err, storage.ErrNotFound) returned false; err=%v", err)
	}
}

func TestErrNotFound_ErrorsIsComparison_Update(t *testing.T) {
	store := storage.NewTaskStore()
	_, err := store.Update(1234, models.UpdateTaskRequest{
		Title:       "Title",
		Description: "Description",
	})
	if err == nil {
		t.Fatal("expected non-nil error for Update with non-existent ID")
	}
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("errors.Is(err, storage.ErrNotFound) returned false; err=%v", err)
	}
}

// --- Compound/integration tests ---

func TestCreate_ThenGetByID_ThenGetMissingID_ReturnsErrNotFound(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.CreateTaskRequest{
		Title:       "Task A",
		Description: "Description A",
	}
	created, err := store.Create(req)
	if err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}

	// Valid lookup must succeed
	_, err = store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID returned error for existing task: %v", err)
	}

	// Lookup with a different ID must return ErrNotFound
	nonExistentID := created.ID + 1000
	_, err = store.GetByID(nonExistentID)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected ErrNotFound for non-existent ID %d, got: %v", nonExistentID, err)
	}
}

func TestCreate_AssignsUniqueIDs(t *testing.T) {
	store := storage.NewTaskStore()
	req := models.CreateTaskRequest{
		Title:       "Task",
		Description: "Description",
	}

	t1, err := store.Create(req)
	if err != nil {
		t.Fatalf("first Create returned unexpected error: %v", err)
	}
	t2, err := store.Create(req)
	if err != nil {
		t.Fatalf("second Create returned unexpected error: %v", err)
	}
	if t1.ID == t2.ID {
		t.Errorf("expected unique IDs for two created tasks, both got ID=%d", t1.ID)
	}
}
