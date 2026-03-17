package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/workspace/repo/internal/handlers"
	"github.com/workspace/repo/internal/models"
	"github.com/workspace/repo/internal/storage"
)

// =============================================================================
// Test Infrastructure
// =============================================================================

// localStore is a minimal in-memory task store used for integration testing.
// It is NOT production code — it exists solely to satisfy handlers.TaskStore
// while the real storage implementation is developed.
type localStore struct {
	mu    sync.Mutex
	tasks map[int64]*models.Task
	idGen *storage.IDGenerator
}

func newLocalStore() *localStore {
	return &localStore{
		tasks: make(map[int64]*models.Task),
		idGen: storage.NewIDGenerator(),
	}
}

// Create implements handlers.TaskStore.
func (s *localStore) Create(req models.CreateTaskRequest) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	task := &models.Task{
		ID:          s.idGen.NextID(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[task.ID] = task
	return task, nil
}

// notImplementedHandler is a placeholder for routes not yet implemented in
// handlers.TaskHandler. Tests targeting these routes will fail until the real
// handler methods are added.
func notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprint(w, `{"error":"not implemented"}`)
}

// newTestServer spins up a real HTTP server backed by a fresh localStore.
// Routes for GET/PUT/DELETE are registered as stubs (501 Not Implemented)
// until the corresponding handlers are added to handlers.TaskHandler.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := newLocalStore()
	h := handlers.NewTaskHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tasks", h.Create)
	// The following routes are not yet implemented in handlers.TaskHandler.
	// Tests targeting these routes are expected to fail until implemented.
	mux.HandleFunc("GET /tasks", notImplementedHandler)
	mux.HandleFunc("GET /tasks/{id}", notImplementedHandler)
	mux.HandleFunc("PUT /tasks/{id}", notImplementedHandler)
	mux.HandleFunc("DELETE /tasks/{id}", notImplementedHandler)

	return httptest.NewServer(mux)
}

// errorResp mirrors the JSON error envelope returned by the API.
type errorResp struct {
	Error string `json:"error"`
}

// doPost sends a POST /tasks request with the given raw JSON body string.
func doPost(t *testing.T, server *httptest.Server, body string) *http.Response {
	t.Helper()
	resp, err := http.Post(server.URL+"/tasks", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST /tasks: %v", err)
	}
	return resp
}

// postTask marshals req to JSON and POSTs it to /tasks, returning the response.
func postTask(t *testing.T, server *httptest.Server, req models.CreateTaskRequest) *http.Response {
	t.Helper()
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return doPost(t, server, string(b))
}

// mustPostTask posts a task and fatally fails if the response is not 201.
// Returns the decoded Task on success.
func mustPostTask(t *testing.T, server *httptest.Server, title, description string) models.Task {
	t.Helper()
	resp := postTask(t, server, models.CreateTaskRequest{Title: title, Description: description})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /tasks: expected 201, got %d", resp.StatusCode)
	}
	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode task: %v", err)
	}
	return task
}

// doRequest sends an HTTP request with the given method, URL, and optional JSON body.
func doRequest(t *testing.T, method, url, body string) *http.Response {
	t.Helper()
	var req *http.Request
	var err error
	if body != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		t.Fatalf("create %s request: %v", method, err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	return resp
}

// =============================================================================
// POST /tasks
// =============================================================================

func TestPost_ValidRequest_Returns201(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp := postTask(t, server, models.CreateTaskRequest{
		Title:       "My Task",
		Description: "A description",
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestPost_ValidRequest_ResponseBodyContainsAllFields(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	req := models.CreateTaskRequest{Title: "Integration Task", Description: "Testing all fields"}
	resp := postTask(t, server, req)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if task.ID <= 0 {
		t.Errorf("expected positive id, got %d", task.ID)
	}
	if task.Title != req.Title {
		t.Errorf("title: expected %q, got %q", req.Title, task.Title)
	}
	if task.Description != req.Description {
		t.Errorf("description: expected %q, got %q", req.Description, task.Description)
	}
	if task.CreatedAt.IsZero() {
		t.Error("created_at should not be zero")
	}
	if task.UpdatedAt.IsZero() {
		t.Error("updated_at should not be zero")
	}
}

func TestPost_ValidRequest_ContentTypeIsJSON(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp := postTask(t, server, models.CreateTaskRequest{
		Title:       "A Task",
		Description: "A description",
	})
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
}

func TestPost_MissingTitle_Returns400WithError(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp := doPost(t, server, `{"description":"some description"}`)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	var e errorResp
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if e.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestPost_EmptyTitle_Returns400WithError(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp := doPost(t, server, `{"title":"","description":"some description"}`)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	var e errorResp
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if e.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestPost_EmptyDescription_Returns400WithError(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp := doPost(t, server, `{"title":"A title","description":""}`)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	var e errorResp
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if e.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestPost_TitleTooLong_Returns400WithError(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	longTitle := strings.Repeat("a", 201)
	body := fmt.Sprintf(`{"title":%q,"description":"some description"}`, longTitle)
	resp := doPost(t, server, body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	var e errorResp
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if e.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestPost_DescriptionTooLong_Returns400WithError(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	longDesc := strings.Repeat("b", 1001)
	body := fmt.Sprintf(`{"title":"A title","description":%q}`, longDesc)
	resp := doPost(t, server, body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
	var e errorResp
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if e.Error == "" {
		t.Error("expected non-empty error message")
	}
}

// =============================================================================
// GET /tasks
// =============================================================================

func TestGetAll_EmptyStore_Returns200WithEmptyArray(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/tasks")
	if err != nil {
		t.Fatalf("GET /tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var tasks []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty array, got %d tasks", len(tasks))
	}
}

func TestGetAll_After3Tasks_Returns200WithAllTasksSortedByID(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	// Create 3 tasks.
	mustPostTask(t, server, "Task One", "First task description")
	mustPostTask(t, server, "Task Two", "Second task description")
	mustPostTask(t, server, "Task Three", "Third task description")

	resp, err := http.Get(server.URL + "/tasks")
	if err != nil {
		t.Fatalf("GET /tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var tasks []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	// Verify sorted by ID ascending.
	if !sort.SliceIsSorted(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	}) {
		t.Error("expected tasks sorted by id ascending")
	}
}

// =============================================================================
// GET /tasks/{id}
// =============================================================================

func TestGetByID_KnownID_Returns200WithTask(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	created := mustPostTask(t, server, "Fetch Me", "A task to fetch by id")

	resp := doRequest(t, http.MethodGet, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if task.ID != created.ID {
		t.Errorf("expected id %d, got %d", created.ID, task.ID)
	}
	if task.Title != created.Title {
		t.Errorf("expected title %q, got %q", created.Title, task.Title)
	}
	if task.Description != created.Description {
		t.Errorf("expected description %q, got %q", created.Description, task.Description)
	}
}

func TestGetByID_UnknownID_Returns404WithError(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp := doRequest(t, http.MethodGet, server.URL+"/tasks/9999", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
	var e errorResp
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if e.Error == "" {
		t.Error("expected non-empty error message")
	}
}

// =============================================================================
// PUT /tasks/{id}
// =============================================================================

func TestPut_ValidUpdate_Returns200WithUpdatedFields(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	created := mustPostTask(t, server, "Original Title", "Original description")

	updateBody := `{"title":"Updated Title","description":"Updated description"}`
	resp := doRequest(t, http.MethodPut, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), updateBody)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var updated models.Task
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if updated.Title != "Updated Title" {
		t.Errorf("title: expected %q, got %q", "Updated Title", updated.Title)
	}
	if updated.Description != "Updated description" {
		t.Errorf("description: expected %q, got %q", "Updated description", updated.Description)
	}
}

func TestPut_ValidUpdate_UpdatedAtIsAfterCreatedAt(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	created := mustPostTask(t, server, "Time Test Task", "Testing updated_at timestamp")

	// Small delay so updated_at > created_at is reliably detectable.
	time.Sleep(10 * time.Millisecond)

	updateBody := `{"title":"Updated Title","description":"Updated description"}`
	resp := doRequest(t, http.MethodPut, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), updateBody)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var updated models.Task
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !updated.UpdatedAt.After(created.CreatedAt) {
		t.Errorf("expected updated_at (%v) to be after created_at (%v)", updated.UpdatedAt, created.CreatedAt)
	}
}

func TestPut_UnknownID_Returns404(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	updateBody := `{"title":"Updated Title","description":"Updated description"}`
	resp := doRequest(t, http.MethodPut, server.URL+"/tasks/9999", updateBody)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestPut_InvalidBody_Returns400(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	created := mustPostTask(t, server, "A Task", "A description")

	resp := doRequest(t, http.MethodPut, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), `{not valid json`)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

// =============================================================================
// DELETE /tasks/{id}
// =============================================================================

func TestDelete_KnownID_Returns204WithNoBody(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	created := mustPostTask(t, server, "Delete Me", "To be deleted")

	resp := doRequest(t, http.MethodDelete, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}

	// Response body should be empty.
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	if buf.Len() != 0 {
		t.Errorf("expected empty body on 204, got %q", buf.String())
	}
}

func TestDelete_UnknownID_Returns404(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	resp := doRequest(t, http.MethodDelete, server.URL+"/tasks/9999", "")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDelete_AfterDelete_GetByIDReturns404(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	created := mustPostTask(t, server, "To Delete", "Will be deleted soon")

	delResp := doRequest(t, http.MethodDelete, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), "")
	delResp.Body.Close()
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE expected 204, got %d", delResp.StatusCode)
	}

	getResp := doRequest(t, http.MethodGet, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), "")
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Errorf("after delete, GET /tasks/%d expected 404, got %d", created.ID, getResp.StatusCode)
	}
}

func TestDelete_AfterDelete_TaskNotIncludedInGetAll(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	created := mustPostTask(t, server, "Ephemeral Task", "Will be removed from list")

	delResp := doRequest(t, http.MethodDelete, fmt.Sprintf("%s/tasks/%d", server.URL, created.ID), "")
	delResp.Body.Close()
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE expected 204, got %d", delResp.StatusCode)
	}

	listResp, err := http.Get(server.URL + "/tasks")
	if err != nil {
		t.Fatalf("GET /tasks: %v", err)
	}
	defer listResp.Body.Close()
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /tasks expected 200, got %d", listResp.StatusCode)
	}

	var tasks []models.Task
	if err := json.NewDecoder(listResp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode tasks: %v", err)
	}
	for _, task := range tasks {
		if task.ID == created.ID {
			t.Errorf("deleted task (id=%d) is still present in GET /tasks", created.ID)
		}
	}
}

// =============================================================================
// ID No-Reuse
// =============================================================================

func TestIDNoReuse_AfterDelete_NewTaskReceivesNewID(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	// Create first task; expect id=1.
	first := mustPostTask(t, server, "First Task", "First task description")
	if first.ID != 1 {
		t.Errorf("expected first task id=1, got %d", first.ID)
	}

	// Delete the first task.
	delResp := doRequest(t, http.MethodDelete, fmt.Sprintf("%s/tasks/%d", server.URL, first.ID), "")
	delResp.Body.Close()
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE expected 204, got %d", delResp.StatusCode)
	}

	// Create a second task — it must NOT reuse id=1.
	second := mustPostTask(t, server, "Second Task", "Second task description")
	if second.ID == first.ID {
		t.Errorf("ID reuse detected: both tasks have id=%d", first.ID)
	}
	if second.ID != 2 {
		t.Errorf("expected second task id=2, got %d", second.ID)
	}
}

// =============================================================================
// Concurrency Safety
// =============================================================================

func TestConcurrency_50GoroutinesPostConcurrently_UniqueIDsNoErrors(t *testing.T) {
	server := newTestServer(t)
	defer server.Close()

	const n = 50
	type result struct {
		id     int64
		status int
	}

	results := make([]result, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			req := models.CreateTaskRequest{
				Title:       fmt.Sprintf("Concurrent Task %d", i+1),
				Description: fmt.Sprintf("Description for concurrent task %d", i+1),
			}
			resp := postTask(t, server, req)
			defer resp.Body.Close()

			results[i].status = resp.StatusCode
			if resp.StatusCode == http.StatusCreated {
				var task models.Task
				if err := json.NewDecoder(resp.Body).Decode(&task); err == nil {
					results[i].id = task.ID
				}
			}
		}()
	}

	wg.Wait()

	// Assert all requests returned 201 (no 500 errors or other failures).
	for i, r := range results {
		if r.status != http.StatusCreated {
			t.Errorf("goroutine %d: expected 201, got %d", i, r.status)
		}
	}

	// Assert all returned IDs are unique and non-zero.
	seen := make(map[int64]bool)
	for i, r := range results {
		if r.id == 0 {
			t.Errorf("goroutine %d: received zero id", i)
			continue
		}
		if seen[r.id] {
			t.Errorf("duplicate id %d detected", r.id)
		}
		seen[r.id] = true
	}
	if len(seen) != n {
		t.Errorf("expected %d unique ids, got %d", n, len(seen))
	}

	// Assert GET /tasks returns all 50 tasks.
	listResp, err := http.Get(server.URL + "/tasks")
	if err != nil {
		t.Fatalf("GET /tasks: %v", err)
	}
	defer listResp.Body.Close()

	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /tasks expected 200, got %d", listResp.StatusCode)
	}

	var tasks []models.Task
	if err := json.NewDecoder(listResp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode tasks: %v", err)
	}
	if len(tasks) != n {
		t.Errorf("expected %d tasks in GET /tasks, got %d", n, len(tasks))
	}
}
