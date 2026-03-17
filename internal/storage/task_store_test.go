package storage_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
)

// newStore is a helper to create a fresh InMemoryTaskStore for each test.
func newStore() *storage.InMemoryTaskStore {
	return storage.NewInMemoryTaskStore()
}

// validCreateReq returns a valid CreateTaskRequest for use in tests.
func validCreateReq(title, description string) models.CreateTaskRequest {
	return models.CreateTaskRequest{
		Title:       title,
		Description: description,
	}
}

// validUpdateReq returns a valid UpdateTaskRequest for use in tests.
func validUpdateReq(title, description string) models.UpdateTaskRequest {
	return models.UpdateTaskRequest{
		Title:       title,
		Description: description,
	}
}

// --- Create tests ---

func TestCreate_ReturnsTaskWithPopulatedID(t *testing.T) {
	s := newStore()
	task, err := s.Create(validCreateReq("Test Task", "Test Description"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if task.ID == 0 {
		t.Error("expected non-zero task ID after create")
	}
}

func TestCreate_ReturnsTaskWithPopulatedCreatedAt(t *testing.T) {
	before := time.Now().UTC()
	s := newStore()
	task, err := s.Create(validCreateReq("Test Task", "Test Description"))
	after := time.Now().UTC()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if task.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt after create")
	}
	if task.CreatedAt.Before(before) || task.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v is not within expected range [%v, %v]", task.CreatedAt, before, after)
	}
}

func TestCreate_ReturnsTaskWithPopulatedUpdatedAt(t *testing.T) {
	before := time.Now().UTC()
	s := newStore()
	task, err := s.Create(validCreateReq("Test Task", "Test Description"))
	after := time.Now().UTC()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if task.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt after create")
	}
	if task.UpdatedAt.Before(before) || task.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt %v is not within expected range [%v, %v]", task.UpdatedAt, before, after)
	}
}

func TestCreate_ReturnsTaskWithCorrectTitle(t *testing.T) {
	s := newStore()
	req := validCreateReq("My Task Title", "My Task Description")
	task, err := s.Create(req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if task.Title != req.Title {
		t.Errorf("expected Title %q, got %q", req.Title, task.Title)
	}
}

func TestCreate_ReturnsTaskWithCorrectDescription(t *testing.T) {
	s := newStore()
	req := validCreateReq("My Task Title", "My Task Description")
	task, err := s.Create(req)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if task.Description != req.Description {
		t.Errorf("expected Description %q, got %q", req.Description, task.Description)
	}
}

func TestCreate_StoredTaskIsRetrievableByID(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Stored Task", "Stored Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	retrieved, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if retrieved.ID != created.ID {
		t.Errorf("ID mismatch: got %d, want %d", retrieved.ID, created.ID)
	}
	if retrieved.Title != created.Title {
		t.Errorf("Title mismatch: got %q, want %q", retrieved.Title, created.Title)
	}
	if retrieved.Description != created.Description {
		t.Errorf("Description mismatch: got %q, want %q", retrieved.Description, created.Description)
	}
}

// --- GetByID tests ---

func TestGetByID_ReturnsTaskForValidID(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Task", "Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	task, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if task.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, task.ID)
	}
}

func TestGetByID_ReturnsErrNotFoundForMissingID(t *testing.T) {
	s := newStore()
	_, err := s.GetByID(99999)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestGetByID_ReturnsCopyNotPointerIntoMap(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Original Title", "Original Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	task, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	// Mutating the returned value should not affect the stored task.
	task.Title = "Mutated Title"

	retrieved, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("second GetByID failed: %v", err)
	}
	if retrieved.Title == "Mutated Title" {
		t.Error("modifying returned task affected the stored task (not a copy)")
	}
}

// --- GetAll tests ---

func TestGetAll_ReturnsEmptySliceForEmptyStore(t *testing.T) {
	s := newStore()
	tasks := s.GetAll()
	if tasks == nil {
		t.Error("expected empty non-nil slice, got nil")
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestGetAll_ReturnsAllTasksAfterMultipleCreates(t *testing.T) {
	s := newStore()

	_, err := s.Create(validCreateReq("Task 1", "Description 1"))
	if err != nil {
		t.Fatalf("create 1 failed: %v", err)
	}
	_, err = s.Create(validCreateReq("Task 2", "Description 2"))
	if err != nil {
		t.Fatalf("create 2 failed: %v", err)
	}
	_, err = s.Create(validCreateReq("Task 3", "Description 3"))
	if err != nil {
		t.Fatalf("create 3 failed: %v", err)
	}

	tasks := s.GetAll()
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestGetAll_ReturnsOneTaskAfterOneCreate(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Only Task", "Only Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	tasks := s.GetAll()
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != created.ID {
		t.Errorf("expected task ID %d, got %d", created.ID, tasks[0].ID)
	}
}

// --- Update tests ---

func TestUpdate_UpdatesTitle(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Original Title", "Original Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	req := validUpdateReq("Updated Title", "Original Description")
	updated, err := s.Update(created.ID, req)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Title != "Updated Title" {
		t.Errorf("expected Title %q, got %q", "Updated Title", updated.Title)
	}
}

func TestUpdate_UpdatesDescription(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Title", "Original Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	req := validUpdateReq("Title", "Updated Description")
	updated, err := s.Update(created.ID, req)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Description != "Updated Description" {
		t.Errorf("expected Description %q, got %q", "Updated Description", updated.Description)
	}
}

func TestUpdate_UpdatesUpdatedAt(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Title", "Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Small sleep to ensure UpdatedAt changes.
	time.Sleep(time.Millisecond)

	before := time.Now().UTC()
	updated, err := s.Update(created.ID, validUpdateReq("New Title", "New Description"))
	after := time.Now().UTC()

	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.UpdatedAt.Before(before) || updated.UpdatedAt.After(after) {
		t.Errorf("UpdatedAt %v is not within expected range [%v, %v]", updated.UpdatedAt, before, after)
	}
}

func TestUpdate_DoesNotChangeCreatedAt(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Title", "Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	originalCreatedAt := created.CreatedAt

	time.Sleep(time.Millisecond)

	updated, err := s.Update(created.ID, validUpdateReq("New Title", "New Description"))
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if !updated.CreatedAt.Equal(originalCreatedAt) {
		t.Errorf("CreatedAt changed after update: got %v, want %v", updated.CreatedAt, originalCreatedAt)
	}
}

func TestUpdate_ReturnsUpdatedTask(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Title", "Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	req := validUpdateReq("New Title", "New Description")
	returned, err := s.Update(created.ID, req)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if returned.Title != "New Title" {
		t.Errorf("expected returned Title %q, got %q", "New Title", returned.Title)
	}
	if returned.Description != "New Description" {
		t.Errorf("expected returned Description %q, got %q", "New Description", returned.Description)
	}
	if returned.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, returned.ID)
	}
}

func TestUpdate_PersistsChangesToStore(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Title", "Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	_, err = s.Update(created.ID, validUpdateReq("Persisted Title", "Persisted Description"))
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	retrieved, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if retrieved.Title != "Persisted Title" {
		t.Errorf("expected Title %q after update, got %q", "Persisted Title", retrieved.Title)
	}
	if retrieved.Description != "Persisted Description" {
		t.Errorf("expected Description %q after update, got %q", "Persisted Description", retrieved.Description)
	}
}

func TestUpdate_ReturnsErrNotFoundForMissingID(t *testing.T) {
	s := newStore()
	_, err := s.Update(99999, validUpdateReq("Title", "Description"))
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

// --- Delete tests ---

func TestDelete_RemovesTask(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Task To Delete", "Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = s.Delete(created.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	tasks := s.GetAll()
	for _, task := range tasks {
		if task.ID == created.ID {
			t.Errorf("task with ID %d still exists after delete", created.ID)
		}
	}
}

func TestDelete_SubsequentGetByIDReturnsErrNotFound(t *testing.T) {
	s := newStore()
	created, err := s.Create(validCreateReq("Task To Delete", "Description"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = s.Delete(created.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}

func TestDelete_ReturnsErrNotFoundForMissingID(t *testing.T) {
	s := newStore()
	err := s.Delete(99999)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestDelete_DoesNotAffectOtherTasks(t *testing.T) {
	s := newStore()
	task1, err := s.Create(validCreateReq("Task 1", "Description 1"))
	if err != nil {
		t.Fatalf("create 1 failed: %v", err)
	}
	task2, err := s.Create(validCreateReq("Task 2", "Description 2"))
	if err != nil {
		t.Fatalf("create 2 failed: %v", err)
	}

	err = s.Delete(task1.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	retrieved, err := s.GetByID(task2.ID)
	if err != nil {
		t.Fatalf("GetByID for task2 failed after deleting task1: %v", err)
	}
	if retrieved.ID != task2.ID {
		t.Errorf("expected task2 ID %d, got %d", task2.ID, retrieved.ID)
	}
}

// --- ID strictly increasing tests ---

func TestCreate_IDsAreStrictlyIncreasing(t *testing.T) {
	s := newStore()

	task1, err := s.Create(validCreateReq("Task 1", "Description 1"))
	if err != nil {
		t.Fatalf("create 1 failed: %v", err)
	}
	task2, err := s.Create(validCreateReq("Task 2", "Description 2"))
	if err != nil {
		t.Fatalf("create 2 failed: %v", err)
	}
	task3, err := s.Create(validCreateReq("Task 3", "Description 3"))
	if err != nil {
		t.Fatalf("create 3 failed: %v", err)
	}

	if task2.ID <= task1.ID {
		t.Errorf("expected ID of task2 (%d) > task1 (%d)", task2.ID, task1.ID)
	}
	if task3.ID <= task2.ID {
		t.Errorf("expected ID of task3 (%d) > task2 (%d)", task3.ID, task2.ID)
	}
}

func TestCreate_IDsAreUniqueAcrossMultipleCreates(t *testing.T) {
	s := newStore()
	seen := make(map[int64]bool)

	for i := 0; i < 10; i++ {
		task, err := s.Create(validCreateReq("Task", "Description"))
		if err != nil {
			t.Fatalf("create %d failed: %v", i, err)
		}
		if seen[task.ID] {
			t.Errorf("duplicate ID %d generated", task.ID)
		}
		seen[task.ID] = true
	}
}

// --- Concurrency safety tests ---

func TestCreate_ConcurrentCreatesSafe(t *testing.T) {
	s := newStore()
	const goroutines = 100

	var wg sync.WaitGroup
	results := make(chan int64, goroutines)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			task, err := s.Create(validCreateReq("Concurrent Task", "Description"))
			if err != nil {
				t.Errorf("concurrent create failed: %v", err)
				return
			}
			results <- task.ID
		}(i)
	}

	wg.Wait()
	close(results)

	// Verify all IDs are unique.
	seen := make(map[int64]bool)
	for id := range results {
		if seen[id] {
			t.Errorf("duplicate ID %d generated under concurrent creates", id)
		}
		seen[id] = true
	}

	if len(seen) != goroutines {
		t.Errorf("expected %d unique IDs, got %d", goroutines, len(seen))
	}
}

func TestGetAll_ConcurrentReadsSafe(t *testing.T) {
	s := newStore()

	// Pre-populate the store.
	for i := 0; i < 10; i++ {
		_, err := s.Create(validCreateReq("Task", "Description"))
		if err != nil {
			t.Fatalf("setup create failed: %v", err)
		}
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			tasks := s.GetAll()
			if len(tasks) != 10 {
				t.Errorf("expected 10 tasks in concurrent GetAll, got %d", len(tasks))
			}
		}()
	}

	wg.Wait()
}

func TestStore_ConcurrentMixedOperationsSafe(t *testing.T) {
	s := newStore()

	// Seed with some initial tasks.
	seed := make([]int64, 5)
	for i := 0; i < 5; i++ {
		task, err := s.Create(validCreateReq("Seed Task", "Description"))
		if err != nil {
			t.Fatalf("seed create failed: %v", err)
		}
		seed[i] = task.ID
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			switch i % 4 {
			case 0:
				s.Create(validCreateReq("Concurrent Task", "Description"))
			case 1:
				s.GetAll()
			case 2:
				s.GetByID(seed[i%5])
			case 3:
				s.Update(seed[i%5], validUpdateReq("Updated Title", "Updated Description"))
			}
		}(i)
	}

	wg.Wait()
}

// --- ErrNotFound sentinel error tests ---

func TestErrNotFound_IsDistinctSentinelError(t *testing.T) {
	if storage.ErrNotFound == nil {
		t.Fatal("ErrNotFound should not be nil")
	}
	if storage.ErrNotFound.Error() == "" {
		t.Fatal("ErrNotFound should have a non-empty message")
	}
}

func TestGetByID_ErrorCanBeIdentifiedWithErrorsIs(t *testing.T) {
	s := newStore()
	_, err := s.GetByID(12345)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected errors.Is(err, ErrNotFound) to be true, got err: %v", err)
	}
}

func TestUpdate_ErrorCanBeIdentifiedWithErrorsIs(t *testing.T) {
	s := newStore()
	_, err := s.Update(12345, validUpdateReq("Title", "Description"))
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected errors.Is(err, ErrNotFound) to be true, got err: %v", err)
	}
}

func TestDelete_ErrorCanBeIdentifiedWithErrorsIs(t *testing.T) {
	s := newStore()
	err := s.Delete(12345)
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("expected errors.Is(err, ErrNotFound) to be true, got err: %v", err)
	}
}

// --- TaskStore interface compliance ---

func TestInMemoryTaskStore_ImplementsTaskStoreInterface(t *testing.T) {
	// This is a compile-time check: if InMemoryTaskStore does not implement
	// TaskStore, this assignment will fail to compile.
	var _ storage.TaskStore = storage.NewInMemoryTaskStore()
}
