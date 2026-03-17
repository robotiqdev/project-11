package storage

import (
	"testing"

	"github.com/workspace/repo/internal/models"
)

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
