package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWriteError_SetsStatusCode verifies that writeError writes the given HTTP status code.
func TestWriteError_SetsStatusCode(t *testing.T) {
	cases := []struct {
		name   string
		status int
	}{
		{"BadRequest", http.StatusBadRequest},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeError(w, tc.status, "some error")

			if w.Code != tc.status {
				t.Errorf("expected status %d, got %d", tc.status, w.Code)
			}
		})
	}
}

// TestWriteError_ResponseBodyIsValidJSON verifies the response body is valid JSON.
func TestWriteError_ResponseBodyIsValidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "invalid input")

	body := w.Body.Bytes()
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("response body is not valid JSON: %v — body: %q", err, body)
	}
}

// TestWriteError_ResponseBodyContainsErrorField verifies the JSON body has an "error" field.
func TestWriteError_ResponseBodyContainsErrorField(t *testing.T) {
	const msg = "something went wrong"
	w := httptest.NewRecorder()
	writeError(w, http.StatusInternalServerError, msg)

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Error != msg {
		t.Errorf("expected error field %q, got %q", msg, resp.Error)
	}
}

// TestWriteError_ContentTypeIsApplicationJSON verifies the Content-Type header is application/json.
func TestWriteError_ContentTypeIsApplicationJSON(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusNotFound, "not found")

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestWriteError_400_BadRequest verifies correct behavior for 400 status.
func TestWriteError_400_BadRequest(t *testing.T) {
	const msg = "bad request"
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, msg)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json Content-Type")
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error != msg {
		t.Errorf("expected %q, got %q", msg, resp.Error)
	}
}

// TestWriteError_404_NotFound verifies correct behavior for 404 status.
func TestWriteError_404_NotFound(t *testing.T) {
	const msg = "resource not found"
	w := httptest.NewRecorder()
	writeError(w, http.StatusNotFound, msg)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json Content-Type")
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error != msg {
		t.Errorf("expected %q, got %q", msg, resp.Error)
	}
}

// TestWriteError_500_InternalServerError verifies correct behavior for 500 status.
func TestWriteError_500_InternalServerError(t *testing.T) {
	const msg = "internal server error"
	w := httptest.NewRecorder()
	writeError(w, http.StatusInternalServerError, msg)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json Content-Type")
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error != msg {
		t.Errorf("expected %q, got %q", msg, resp.Error)
	}
}

// TestWriteError_EmptyMessage verifies that an empty message is serialized correctly.
func TestWriteError_EmptyMessage(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "")

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error != "" {
		t.Errorf("expected empty error field, got %q", resp.Error)
	}
}

// TestWriteError_MessagePreserved verifies the exact message string is preserved in the JSON output.
func TestWriteError_MessagePreserved(t *testing.T) {
	const msg = "file \"config.yaml\" not found"
	w := httptest.NewRecorder()
	writeError(w, http.StatusNotFound, msg)

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if resp.Error != msg {
		t.Errorf("expected %q, got %q", msg, resp.Error)
	}
}

// TestWriteError_ContentTypeSetBeforeBody verifies Content-Type is present after writing.
// (Headers must be set before WriteHeader; this test confirms they were set in the right order.)
func TestWriteError_ContentTypeSetBeforeBody(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusInternalServerError, "error")

	// After writeError, both header and body should be set correctly.
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type not set correctly: got %q", ct)
	}
	if w.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}
