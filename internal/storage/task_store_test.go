package storage

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/workspace/repo/internal/models"
)

// newStore is a helper to create a fresh InMemoryTaskStore for each test.
func newStore() *InMemoryTaskStore {
	return NewInMemoryTaskStore()
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
	if !errors.Is(err, ErrNotFound) {
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
	if !errors.Is(err, ErrNotFound) {
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
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}

func TestDelete_ReturnsErrNotFoundForMissingID(t *testing.T) {
	s := newStore()
	err := s.Delete(99999)
	if err == nil {
		t.Fatal("expected error for missing ID, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
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
	if ErrNotFound == nil {
		t.Fatal("ErrNotFound should not be nil")
	}
	if ErrNotFound.Error() == "" {
		t.Fatal("ErrNotFound should have a non-empty message")
	}
}

func TestGetByID_ErrorCanBeIdentifiedWithErrorsIs(t *testing.T) {
	s := newStore()
	_, err := s.GetByID(12345)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected errors.Is(err, ErrNotFound) to be true, got err: %v", err)
	}
}

func TestUpdate_ErrorCanBeIdentifiedWithErrorsIs(t *testing.T) {
	s := newStore()
	_, err := s.Update(12345, validUpdateReq("Title", "Description"))
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected errors.Is(err, ErrNotFound) to be true, got err: %v", err)
	}
}

func TestDelete_ErrorCanBeIdentifiedWithErrorsIs(t *testing.T) {
	s := newStore()
	err := s.Delete(12345)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected errors.Is(err, ErrNotFound) to be true, got err: %v", err)
	}
}

// --- TaskStore interface compliance ---

func TestInMemoryTaskStore_ImplementsTaskStoreInterface(t *testing.T) {
	// This is a compile-time check: if InMemoryTaskStore does not implement
	// TaskStore, this assignment will fail to compile.
	var _ TaskStore = NewInMemoryTaskStore()
}

// --- NewInMemoryTaskStore constructor tests ---

// TestNewInMemoryTaskStore_ReturnsNonNilStore verifies the constructor returns a non-nil pointer.
func TestNewInMemoryTaskStore_ReturnsNonNilStore(t *testing.T) {
	store := NewInMemoryTaskStore()
	if store == nil {
		t.Fatal("NewInMemoryTaskStore() returned nil; expected a valid *InMemoryTaskStore")
	}
}

// TestNewInMemoryTaskStore_TasksMapIsNotNil verifies the tasks map is initialized (not nil)
// to prevent nil map assignment panics on write operations.
func TestNewInMemoryTaskStore_TasksMapIsNotNil(t *testing.T) {
	store := NewInMemoryTaskStore()
	if store.tasks == nil {
		t.Fatal("NewInMemoryTaskStore() returned store with nil tasks map; map must be initialized via make()")
	}
}

// TestNewInMemoryTaskStore_TasksMapIsEmpty verifies the tasks map starts with zero entries.
func TestNewInMemoryTaskStore_TasksMapIsEmpty(t *testing.T) {
	store := NewInMemoryTaskStore()
	if len(store.tasks) != 0 {
		t.Fatalf("NewInMemoryTaskStore() returned store with %d tasks; expected 0", len(store.tasks))
	}
}

// TestNewInMemoryTaskStore_IDGeneratorIsNotNil verifies the IDGenerator is initialized.
func TestNewInMemoryTaskStore_IDGeneratorIsNotNil(t *testing.T) {
	store := NewInMemoryTaskStore()
	if store.idGen == nil {
		t.Fatal("NewInMemoryTaskStore() returned store with nil idGen; IDGenerator must be initialized via NewIDGenerator()")
	}
}

// TestNewInMemoryTaskStore_IsUsableImmediately verifies that the store can be used
// for write operations immediately after construction without causing a panic.
// A nil map would panic on assignment.
func TestNewInMemoryTaskStore_IsUsableImmediately(t *testing.T) {
	store := NewInMemoryTaskStore()

	// Writing to a nil map panics; this test ensures the map is safe to write to.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked when writing to store.tasks immediately after NewInMemoryTaskStore(): %v", r)
		}
	}()

	store.tasks[1] = models.Task{ID: 1, Title: "Test", Description: "Test description"}
}

// TestNewInMemoryTaskStore_TasksMapSupportsReadAfterWrite verifies that tasks written
// to the map can be read back correctly.
func TestNewInMemoryTaskStore_TasksMapSupportsReadAfterWrite(t *testing.T) {
	store := NewInMemoryTaskStore()

	if store.tasks == nil {
		t.Fatal("cannot test read-after-write: store.tasks is nil (map not initialized)")
	}

	task := models.Task{ID: 42, Title: "My Task", Description: "My Description"}
	store.tasks[42] = task

	got, ok := store.tasks[42]
	if !ok {
		t.Fatal("expected to find task with ID 42 after writing it, but it was not found")
	}
	if got.ID != task.ID {
		t.Errorf("task ID mismatch: got %d, want %d", got.ID, task.ID)
	}
	if got.Title != task.Title {
		t.Errorf("task Title mismatch: got %q, want %q", got.Title, task.Title)
	}
}

// TestNewInMemoryTaskStore_MultipleCallsReturnIndependentStores verifies that each
// call to NewInMemoryTaskStore returns a store with its own independent map,
// so mutations in one store do not affect another.
func TestNewInMemoryTaskStore_MultipleCallsReturnIndependentStores(t *testing.T) {
	store1 := NewInMemoryTaskStore()
	store2 := NewInMemoryTaskStore()

	if store1.tasks == nil || store2.tasks == nil {
		t.Fatal("cannot test independence: one or both stores have nil tasks maps")
	}

	store1.tasks[1] = models.Task{ID: 1, Title: "Store1 Task", Description: "desc"}

	if _, ok := store2.tasks[1]; ok {
		t.Fatal("writing to store1 should not affect store2; stores must have independent maps")
	}
}

// TestInMemoryTaskStore_ZeroValueHasNilMap documents that a zero-value InMemoryTaskStore
// has a nil tasks map, which is unsafe for writes. This confirms why NewInMemoryTaskStore()
// must explicitly initialize the map.
func TestInMemoryTaskStore_ZeroValueHasNilMap(t *testing.T) {
	var store InMemoryTaskStore
	if store.tasks != nil {
		t.Fatal("zero-value InMemoryTaskStore unexpectedly has a non-nil tasks map; this documents expected zero-value behavior")
	}
}

// TestInMemoryTaskStore_ZeroValueHasNilIDGenerator documents that a zero-value
// InMemoryTaskStore has a nil idGen, confirming why NewInMemoryTaskStore() must
// explicitly initialize it.
func TestInMemoryTaskStore_ZeroValueHasNilIDGenerator(t *testing.T) {
	var store InMemoryTaskStore
	if store.idGen != nil {
		t.Fatal("zero-value InMemoryTaskStore unexpectedly has a non-nil idGen; this documents expected zero-value behavior")
	}
}

// TestNewInMemoryTaskStore_WritingToZeroValuePanics documents that writing to a
// zero-value store's nil map panics, reinforcing the need for the constructor.
func TestNewInMemoryTaskStore_WritingToZeroValuePanics(t *testing.T) {
	var store InMemoryTaskStore

	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		// This should panic because store.tasks is nil.
		store.tasks[1] = models.Task{ID: 1, Title: "Test", Description: "desc"}
	}()

	if !panicked {
		t.Fatal("expected panic when writing to zero-value store's nil map, but no panic occurred")
	}
}

// TestNewInMemoryTaskStore_StructHasExpectedFields verifies the InMemoryTaskStore
// struct has the expected exported-accessible fields/types at compile time.
// This is a compile-time assertion via type usage.
func TestNewInMemoryTaskStore_StructHasExpectedFields(t *testing.T) {
	store := NewInMemoryTaskStore()

	// Verify tasks field is map[int64]models.Task by assigning a typed value.
	var _ map[int64]models.Task = store.tasks

	// Verify idGen field is *IDGenerator by assigning its type.
	var _ *IDGenerator = store.idGen

	// If the code compiles, the struct has the correct field types.
	_ = store
}

// ── Create ────────────────────────────────────────────────────────────────────

// TestCreate_ValidRequest_ReturnsTaskWithPopulatedFields verifies that a valid
// CreateTaskRequest produces a Task with a non-zero ID, the correct Title and
// Description, and non-zero CreatedAt / UpdatedAt timestamps.
func TestCreate_ValidRequest_ReturnsTaskWithPopulatedFields(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	req := models.CreateTaskRequest{Title: "My Task", Description: "Some description"}

	task, err := store.Create(req)
	if err != nil {
		t.Fatalf("Create() returned unexpected error: %v", err)
	}
	if task.ID == 0 {
		t.Error("Create() returned task with ID 0; expected non-zero ID")
	}
	if task.Title != req.Title {
		t.Errorf("Create() title = %q; want %q", task.Title, req.Title)
	}
	if task.Description != req.Description {
		t.Errorf("Create() description = %q; want %q", task.Description, req.Description)
	}
	if task.CreatedAt.IsZero() {
		t.Error("Create() returned task with zero CreatedAt; expected a timestamp")
	}
	if task.UpdatedAt.IsZero() {
		t.Error("Create() returned task with zero UpdatedAt; expected a timestamp")
	}
}

// TestCreate_ValidRequest_TaskStoredAndRetrievableViaGetByID confirms that a
// successfully created task can be fetched back by its ID.
func TestCreate_ValidRequest_TaskStoredAndRetrievableViaGetByID(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	req := models.CreateTaskRequest{Title: "Stored Task", Description: "desc"}

	created, err := store.Create(req)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	got, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID(%d) error after Create: %v", created.ID, err)
	}
	if got.ID != created.ID {
		t.Errorf("GetByID returned ID %d; want %d", got.ID, created.ID)
	}
}

// TestCreate_MultipleRequests_ReturnDistinctIDs ensures successive Create calls
// produce tasks with unique, distinct IDs.
func TestCreate_MultipleRequests_ReturnDistinctIDs(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()

	type tc struct {
		name string
		req  models.CreateTaskRequest
	}
	cases := []tc{
		{"first", models.CreateTaskRequest{Title: "Task 1", Description: "desc 1"}},
		{"second", models.CreateTaskRequest{Title: "Task 2", Description: "desc 2"}},
		{"third", models.CreateTaskRequest{Title: "Task 3", Description: "desc 3"}},
	}

	ids := make(map[int64]bool)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			task, err := store.Create(c.req)
			if err != nil {
				t.Fatalf("Create() error: %v", err)
			}
			if ids[task.ID] {
				t.Errorf("Create() returned duplicate ID %d", task.ID)
			}
			ids[task.ID] = true
		})
	}
}

// TestCreate_EmptyTitle_ReturnsError verifies that Create rejects a request with
// an empty title and returns a non-nil error.
func TestCreate_EmptyTitle_ReturnsError(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	req := models.CreateTaskRequest{Title: "", Description: "some description"}

	_, err := store.Create(req)
	if err == nil {
		t.Fatal("Create() with empty title returned nil error; expected validation error")
	}
}

// TestCreate_EmptyDescription_ReturnsError verifies that Create rejects a request
// with an empty description and returns a non-nil error.
func TestCreate_EmptyDescription_ReturnsError(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	req := models.CreateTaskRequest{Title: "Valid Title", Description: ""}

	_, err := store.Create(req)
	if err == nil {
		t.Fatal("Create() with empty description returned nil error; expected validation error")
	}
}

// TestCreate_ErrorCases_TableDriven covers multiple invalid request inputs.
func TestCreate_ErrorCases_TableDriven(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		req  models.CreateTaskRequest
	}{
		{"empty_title", models.CreateTaskRequest{Title: "", Description: "desc"}},
		{"empty_description", models.CreateTaskRequest{Title: "title", Description: ""}},
		{"both_empty", models.CreateTaskRequest{Title: "", Description: ""}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			store := NewInMemoryTaskStore()
			_, err := store.Create(tc.req)
			if err == nil {
				t.Errorf("Create(%+v) returned nil error; expected error", tc.req)
			}
		})
	}
}

// ── GetByID ───────────────────────────────────────────────────────────────────

// TestGetByID_ExistingID_ReturnsTask verifies that GetByID returns the correct
// Task when the ID exists in the store.
func TestGetByID_ExistingID_ReturnsTask(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	req := models.CreateTaskRequest{Title: "Find Me", Description: "I exist"}
	created, err := store.Create(req)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	got, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID(%d) returned error: %v", created.ID, err)
	}
	if got.ID != created.ID {
		t.Errorf("GetByID ID = %d; want %d", got.ID, created.ID)
	}
	if got.Title != req.Title {
		t.Errorf("GetByID Title = %q; want %q", got.Title, req.Title)
	}
	if got.Description != req.Description {
		t.Errorf("GetByID Description = %q; want %q", got.Description, req.Description)
	}
}

// TestGetByID_NonExistentID_ReturnsErrNotFound verifies that GetByID returns
// ErrNotFound when no task with that ID exists.
func TestGetByID_NonExistentID_ReturnsErrNotFound(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()

	_, err := store.GetByID(9999)
	if err == nil {
		t.Fatal("GetByID(9999) returned nil error for non-existent ID; expected ErrNotFound")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("GetByID error = %v; want ErrNotFound", err)
	}
}

// TestGetByID_TableDriven covers found and not-found lookup scenarios.
func TestGetByID_TableDriven(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "Table Task", Description: "desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	cases := []struct {
		name        string
		id          int64
		wantErr     bool
		wantErrType error
	}{
		{"found", created.ID, false, nil},
		{"not_found_zero", 0, true, ErrNotFound},
		{"not_found_large", 99999, true, ErrNotFound},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			task, err := store.GetByID(tc.id)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetByID(%d) returned nil error; want error", tc.id)
				}
				if tc.wantErrType != nil && !errors.Is(err, tc.wantErrType) {
					t.Errorf("GetByID(%d) error = %v; want %v", tc.id, err, tc.wantErrType)
				}
			} else {
				if err != nil {
					t.Fatalf("GetByID(%d) unexpected error: %v", tc.id, err)
				}
				if task.ID != tc.id {
					t.Errorf("GetByID(%d) returned task.ID = %d; want %d", tc.id, task.ID, tc.id)
				}
			}
		})
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

// TestUpdate_ExistingID_ReturnsUpdatedTask verifies that Update returns a Task
// with the new field values when the ID exists.
func TestUpdate_ExistingID_ReturnsUpdatedTask(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "Old Title", Description: "Old desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	updateReq := models.UpdateTaskRequest{Title: "New Title", Description: "New desc"}
	updated, err := store.Update(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Update(%d) returned error: %v", created.ID, err)
	}
	if updated.Title != updateReq.Title {
		t.Errorf("Update Title = %q; want %q", updated.Title, updateReq.Title)
	}
	if updated.Description != updateReq.Description {
		t.Errorf("Update Description = %q; want %q", updated.Description, updateReq.Description)
	}
	if updated.ID != created.ID {
		t.Errorf("Update returned task with ID %d; want %d", updated.ID, created.ID)
	}
}

// TestUpdate_ExistingID_FieldsActuallyChangedInStore confirms that after Update,
// GetByID reflects the new values (not the old ones).
func TestUpdate_ExistingID_FieldsActuallyChangedInStore(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "Original", Description: "Original desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	_, err = store.Update(created.ID, models.UpdateTaskRequest{Title: "Changed", Description: "Changed desc"})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	got, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID after Update error: %v", err)
	}
	if got.Title != "Changed" {
		t.Errorf("GetByID after Update Title = %q; want %q", got.Title, "Changed")
	}
	if got.Description != "Changed desc" {
		t.Errorf("GetByID after Update Description = %q; want %q", got.Description, "Changed desc")
	}
}

// TestUpdate_ExistingID_UpdatedAtIsRefreshed verifies that UpdatedAt is advanced
// after an update (it must not remain the same as CreatedAt if time has passed).
func TestUpdate_ExistingID_UpdatedAtIsRefreshed(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "Task", Description: "desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	// Sleep briefly so clock advances (only needed if time resolution allows).
	time.Sleep(time.Millisecond)

	updated, err := store.Update(created.ID, models.UpdateTaskRequest{Title: "Task", Description: "updated"})
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}
	if updated.UpdatedAt.IsZero() {
		t.Error("Update returned task with zero UpdatedAt; expected a timestamp")
	}
	// UpdatedAt must not be before CreatedAt.
	if updated.UpdatedAt.Before(created.CreatedAt) {
		t.Errorf("UpdatedAt %v is before CreatedAt %v; expected UpdatedAt >= CreatedAt",
			updated.UpdatedAt, created.CreatedAt)
	}
}

// TestUpdate_NonExistentID_ReturnsErrNotFound verifies that Update returns
// ErrNotFound when no task with that ID exists.
func TestUpdate_NonExistentID_ReturnsErrNotFound(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	req := models.UpdateTaskRequest{Title: "Title", Description: "desc"}

	_, err := store.Update(9999, req)
	if err == nil {
		t.Fatal("Update(9999) returned nil error for non-existent ID; expected ErrNotFound")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Update error = %v; want ErrNotFound", err)
	}
}

// TestUpdate_TableDriven covers multiple found/not-found scenarios.
func TestUpdate_TableDriven(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "Existing", Description: "desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	cases := []struct {
		name        string
		id          int64
		req         models.UpdateTaskRequest
		wantErr     bool
		wantErrType error
		wantTitle   string
		wantDesc    string
	}{
		{
			name:      "found_changes_fields",
			id:        created.ID,
			req:       models.UpdateTaskRequest{Title: "Updated", Description: "Updated desc"},
			wantErr:   false,
			wantTitle: "Updated",
			wantDesc:  "Updated desc",
		},
		{
			name:        "not_found",
			id:          88888,
			req:         models.UpdateTaskRequest{Title: "X", Description: "Y"},
			wantErr:     true,
			wantErrType: ErrNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			task, err := store.Update(tc.id, tc.req)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Update(%d) returned nil error; want error", tc.id)
				}
				if tc.wantErrType != nil && !errors.Is(err, tc.wantErrType) {
					t.Errorf("Update(%d) error = %v; want %v", tc.id, err, tc.wantErrType)
				}
			} else {
				if err != nil {
					t.Fatalf("Update(%d) unexpected error: %v", tc.id, err)
				}
				if task.Title != tc.wantTitle {
					t.Errorf("Update Title = %q; want %q", task.Title, tc.wantTitle)
				}
				if task.Description != tc.wantDesc {
					t.Errorf("Update Description = %q; want %q", task.Description, tc.wantDesc)
				}
			}
		})
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

// TestDelete_ExistingID_RemovesTaskFromStore verifies that Delete succeeds for an
// existing task and that the task is no longer retrievable via GetByID.
func TestDelete_ExistingID_RemovesTaskFromStore(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "To Delete", Description: "desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	if err := store.Delete(created.ID); err != nil {
		t.Fatalf("Delete(%d) returned unexpected error: %v", created.ID, err)
	}

	_, err = store.GetByID(created.ID)
	if err == nil {
		t.Errorf("GetByID(%d) returned nil error after Delete; expected ErrNotFound", created.ID)
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("GetByID after Delete error = %v; want ErrNotFound", err)
	}
}

// TestDelete_NonExistentID_ReturnsErrNotFound verifies that deleting a task that
// does not exist returns ErrNotFound.
func TestDelete_NonExistentID_ReturnsErrNotFound(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()

	err := store.Delete(9999)
	if err == nil {
		t.Fatal("Delete(9999) returned nil error for non-existent ID; expected ErrNotFound")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete error = %v; want ErrNotFound", err)
	}
}

// TestDelete_IdempotencyBehavior_SecondDeleteReturnsErrNotFound verifies that
// deleting an already-deleted task is treated as not-found (not silent success).
func TestDelete_IdempotencyBehavior_SecondDeleteReturnsErrNotFound(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "Delete Twice", Description: "desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	if err := store.Delete(created.ID); err != nil {
		t.Fatalf("first Delete(%d) error: %v", created.ID, err)
	}

	err = store.Delete(created.ID)
	if err == nil {
		t.Fatal("second Delete of same ID returned nil error; expected ErrNotFound")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("second Delete error = %v; want ErrNotFound", err)
	}
}

// TestDelete_TableDriven covers found and not-found deletion scenarios.
func TestDelete_TableDriven(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		setup       func(s *InMemoryTaskStore) int64
		wantErr     bool
		wantErrType error
	}{
		{
			name: "existing_id_succeeds",
			setup: func(s *InMemoryTaskStore) int64 {
				task, _ := s.Create(models.CreateTaskRequest{Title: "T", Description: "D"})
				return task.ID
			},
			wantErr: false,
		},
		{
			name: "non_existent_id_returns_err_not_found",
			setup: func(_ *InMemoryTaskStore) int64 {
				return 77777
			},
			wantErr:     true,
			wantErrType: ErrNotFound,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			store := NewInMemoryTaskStore()
			id := tc.setup(store)

			err := store.Delete(id)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Delete(%d) returned nil error; want error", id)
				}
				if tc.wantErrType != nil && !errors.Is(err, tc.wantErrType) {
					t.Errorf("Delete(%d) error = %v; want %v", id, err, tc.wantErrType)
				}
			} else {
				if err != nil {
					t.Fatalf("Delete(%d) unexpected error: %v", id, err)
				}
			}
		})
	}
}

// ── GetAll ────────────────────────────────────────────────────────────────────

// TestGetAll_EmptyStore_ReturnsEmptySlice verifies that GetAll on an empty store
// returns an empty (or nil) slice with no error.
func TestGetAll_EmptyStore_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	tasks := store.GetAll()
	if len(tasks) != 0 {
		t.Errorf("GetAll() on empty store returned %d tasks; want 0", len(tasks))
	}
}

// TestGetAll_AfterCreates_ReturnsAllTasks verifies that GetAll returns a slice
// containing every task that was created.
func TestGetAll_AfterCreates_ReturnsAllTasks(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	reqs := []models.CreateTaskRequest{
		{Title: "Task A", Description: "desc A"},
		{Title: "Task B", Description: "desc B"},
		{Title: "Task C", Description: "desc C"},
	}

	createdIDs := make(map[int64]bool)
	for _, r := range reqs {
		task, err := store.Create(r)
		if err != nil {
			t.Fatalf("Create() error: %v", err)
		}
		createdIDs[task.ID] = true
	}

	all := store.GetAll()
	if len(all) != len(reqs) {
		t.Errorf("GetAll() returned %d tasks; want %d", len(all), len(reqs))
	}
	for _, task := range all {
		if !createdIDs[task.ID] {
			t.Errorf("GetAll() returned unexpected task with ID %d", task.ID)
		}
	}
}

// TestGetAll_AfterDelete_DoesNotIncludeDeletedTask verifies that a deleted task
// no longer appears in the GetAll result.
func TestGetAll_AfterDelete_DoesNotIncludeDeletedTask(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	t1, _ := store.Create(models.CreateTaskRequest{Title: "Keep", Description: "desc"})
	t2, _ := store.Create(models.CreateTaskRequest{Title: "Delete Me", Description: "desc"})

	if err := store.Delete(t2.ID); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	all := store.GetAll()
	for _, task := range all {
		if task.ID == t2.ID {
			t.Errorf("GetAll() still contains deleted task ID %d", t2.ID)
		}
	}
	_ = t1
}

// TestGetAll_ReturnsCopies_MutatingSliceDoesNotAffectStore verifies that the
// Tasks returned by GetAll are independent copies; mutating them does not corrupt
// the store's internal state.
func TestGetAll_ReturnsCopies_MutatingSliceDoesNotAffectStore(t *testing.T) {
	t.Parallel()

	store := NewInMemoryTaskStore()
	created, err := store.Create(models.CreateTaskRequest{Title: "Original", Description: "desc"})
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	all := store.GetAll()
	if len(all) != 1 {
		t.Fatalf("GetAll() returned %d tasks; want 1", len(all))
	}

	// Mutate the returned copy.
	all[0].Title = "Mutated"

	// The store should still have the original.
	got, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}
	if got.Title != "Original" {
		t.Errorf("Store title = %q after mutating GetAll result; want %q", got.Title, "Original")
	}
}

