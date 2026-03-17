package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/workspace/repo/internal/router"
	"github.com/workspace/repo/internal/storage"
)

// newTestHandler returns an http.Handler backed by an empty in-memory store.
func newTestHandler() http.Handler {
	store := storage.NewInMemoryTaskStore()
	return router.New(store)
}

// newTestServer starts a real httptest.Server using the router under test.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(newTestHandler())
	t.Cleanup(srv.Close)
	return srv
}

// ---------------------------------------------------------------------------
// 404 smoke tests — unknown routes must not be served
// ---------------------------------------------------------------------------

func TestRouter_UnknownRoute_Returns404(t *testing.T) {
	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/unknown")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /unknown: expected 404, got %d", resp.StatusCode)
	}
}

func TestRouter_DeepUnknownRoute_Returns404(t *testing.T) {
	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/api/v1/tasks")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /api/v1/tasks: expected 404, got %d", resp.StatusCode)
	}
}

func TestRouter_RootPath_Returns404(t *testing.T) {
	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /: expected 404, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// 405 tests — registered paths called with wrong HTTP methods must return 405
// ---------------------------------------------------------------------------

func TestRouter_PatchTasks_Returns405(t *testing.T) {
	srv := newTestServer(t)

	req, err := http.NewRequest(http.MethodPatch, srv.URL+"/tasks", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("PATCH /tasks: expected 405, got %d", resp.StatusCode)
	}
}

func TestRouter_DeleteTasks_Returns405(t *testing.T) {
	// DELETE is not registered on the /tasks collection, only on /tasks/{id}.
	srv := newTestServer(t)

	req, err := http.NewRequest(http.MethodDelete, srv.URL+"/tasks", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("DELETE /tasks: expected 405, got %d", resp.StatusCode)
	}
}

func TestRouter_PutTasks_Returns405(t *testing.T) {
	// PUT is not registered on the /tasks collection, only on /tasks/{id}.
	srv := newTestServer(t)

	req, err := http.NewRequest(http.MethodPut, srv.URL+"/tasks", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("PUT /tasks: expected 405, got %d", resp.StatusCode)
	}
}

func TestRouter_PatchTaskByID_Returns405(t *testing.T) {
	srv := newTestServer(t)

	req, err := http.NewRequest(http.MethodPatch, srv.URL+"/tasks/1", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("PATCH /tasks/1: expected 405, got %d", resp.StatusCode)
	}
}

func TestRouter_PostTaskByID_Returns405(t *testing.T) {
	// POST is only registered on the collection, not on /tasks/{id}.
	srv := newTestServer(t)

	req, err := http.NewRequest(http.MethodPost, srv.URL+"/tasks/1", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("POST /tasks/1: expected 405, got %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// Smoke tests — registered routes must be reachable (2xx or known 4xx)
// ---------------------------------------------------------------------------

func TestRouter_GetTasks_IsReachable(t *testing.T) {
	// GET /tasks must be served by the router (not 404/405).
	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/tasks")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		t.Errorf("GET /tasks: expected route to be registered, got %d", resp.StatusCode)
	}
}

func TestRouter_GetTaskByID_IsReachable(t *testing.T) {
	// GET /tasks/{id} must be served by the router (not 404 from router itself).
	// A 404 returned by the handler (task not found) is acceptable.
	srv := newTestServer(t)

	resp, err := http.Get(srv.URL + "/tasks/999")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMethodNotAllowed {
		t.Errorf("GET /tasks/999: expected route to be registered, got 405")
	}
}
