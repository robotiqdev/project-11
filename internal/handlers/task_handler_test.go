package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/workspace/repo/internal/models"
)

// mockTaskStore is a test double for TaskStore.
type mockTaskStore struct {
	createFn func(req models.CreateTaskRequest) (*models.Task, error)
}

func (m *mockTaskStore) Create(req models.CreateTaskRequest) (*models.Task, error) {
	return m.createFn(req)
}

// newSuccessStore returns a mockTaskStore that always creates a task successfully.
func newSuccessStore() *mockTaskStore {
	return &mockTaskStore{
		createFn: func(req models.CreateTaskRequest) (*models.Task, error) {
			now := time.Now()
			return &models.Task{
				ID:          1,
				Title:       req.Title,
				Description: req.Description,
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}
}

// postJSON sends a POST request with a JSON body to the handler's Create method.
func postJSON(h *TaskHandler, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)
	return w
}

// --- Tests ---

// TestCreate_ValidRequest_Returns201 verifies that a valid POST returns HTTP 201.
func TestCreate_ValidRequest_Returns201(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Buy milk","description":"Get 2 litres of full-fat milk"}`
	w := postJSON(h, body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

// TestCreate_ValidRequest_ResponseContainsID verifies the response body includes a positive id.
func TestCreate_ValidRequest_ResponseContainsID(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Buy milk","description":"Get some milk"}`
	w := postJSON(h, body)

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID <= 0 {
		t.Errorf("expected id > 0, got %d", resp.ID)
	}
}

// TestCreate_ValidRequest_ResponseContainsTitle verifies the response body includes the title.
func TestCreate_ValidRequest_ResponseContainsTitle(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	const title = "Buy milk"
	body := fmt.Sprintf(`{"title":%q,"description":"Get some milk"}`, title)
	w := postJSON(h, body)

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Title != title {
		t.Errorf("expected title %q, got %q", title, resp.Title)
	}
}

// TestCreate_ValidRequest_ResponseContainsDescription verifies the response body includes the description.
func TestCreate_ValidRequest_ResponseContainsDescription(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	const desc = "Get 2 litres of full-fat milk"
	body := fmt.Sprintf(`{"title":"Buy milk","description":%q}`, desc)
	w := postJSON(h, body)

	var resp models.Task
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Description != desc {
		t.Errorf("expected description %q, got %q", desc, resp.Description)
	}
}

// TestCreate_ValidRequest_ResponseContainsCreatedAt verifies the response includes created_at.
func TestCreate_ValidRequest_ResponseContainsCreatedAt(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Buy milk","description":"Get some milk"}`
	w := postJSON(h, body)

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := raw["created_at"]; !ok {
		t.Error("expected created_at field in response")
	}
}

// TestCreate_ValidRequest_ResponseContainsUpdatedAt verifies the response includes updated_at.
func TestCreate_ValidRequest_ResponseContainsUpdatedAt(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Buy milk","description":"Get some milk"}`
	w := postJSON(h, body)

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := raw["updated_at"]; !ok {
		t.Error("expected updated_at field in response")
	}
}

// TestCreate_ValidRequest_ContentTypeIsJSON verifies the response Content-Type is application/json.
func TestCreate_ValidRequest_ContentTypeIsJSON(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Buy milk","description":"Get some milk"}`
	w := postJSON(h, body)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}
}

// TestCreate_MissingTitle_Returns400 verifies that omitting title returns HTTP 400.
func TestCreate_MissingTitle_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"description":"A task without a title"}`
	w := postJSON(h, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestCreate_MissingTitle_ErrorMessageIsTitleRequired verifies the error body says "title is required".
func TestCreate_MissingTitle_ErrorMessageIsTitleRequired(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"description":"A task without a title"}`
	w := postJSON(h, body)

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	const want = "title is required"
	if resp.Error != want {
		t.Errorf("expected error %q, got %q", want, resp.Error)
	}
}

// TestCreate_EmptyTitle_Returns400 verifies that an empty title string returns HTTP 400.
func TestCreate_EmptyTitle_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"","description":"A task with empty title"}`
	w := postJSON(h, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty title, got %d", w.Code)
	}
}

// TestCreate_TitleExceeds200Chars_Returns400 verifies that a title longer than 200 chars returns 400.
func TestCreate_TitleExceeds200Chars_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	longTitle := strings.Repeat("a", 201)
	body := fmt.Sprintf(`{"title":%q,"description":"Some description"}`, longTitle)
	w := postJSON(h, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for title > 200 chars, got %d", w.Code)
	}
}

// TestCreate_TitleExactly200Chars_Returns201 verifies that a title of exactly 200 chars is accepted.
func TestCreate_TitleExactly200Chars_Returns201(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	exactTitle := strings.Repeat("a", 200)
	body := fmt.Sprintf(`{"title":%q,"description":"Some description"}`, exactTitle)
	w := postJSON(h, body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 for title of exactly 200 chars, got %d", w.Code)
	}
}

// TestCreate_EmptyDescription_Returns400 verifies that an empty description returns 400.
func TestCreate_EmptyDescription_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Valid Title","description":""}`
	w := postJSON(h, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty description, got %d", w.Code)
	}
}

// TestCreate_MissingDescription_Returns400 verifies that omitting description returns 400.
func TestCreate_MissingDescription_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Valid Title"}`
	w := postJSON(h, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing description, got %d", w.Code)
	}
}

// TestCreate_DescriptionExceeds1000Chars_Returns400 verifies that a description > 1000 chars returns 400.
func TestCreate_DescriptionExceeds1000Chars_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	longDesc := strings.Repeat("d", 1001)
	body := fmt.Sprintf(`{"title":"Valid Title","description":%q}`, longDesc)
	w := postJSON(h, body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for description > 1000 chars, got %d", w.Code)
	}
}

// TestCreate_DescriptionExactly1000Chars_Returns201 verifies that a description of exactly 1000 chars is accepted.
func TestCreate_DescriptionExactly1000Chars_Returns201(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	exactDesc := strings.Repeat("d", 1000)
	body := fmt.Sprintf(`{"title":"Valid Title","description":%q}`, exactDesc)
	w := postJSON(h, body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201 for description of exactly 1000 chars, got %d", w.Code)
	}
}

// TestCreate_InvalidJSON_Returns400 verifies that malformed JSON returns 400.
func TestCreate_InvalidJSON_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{not valid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d", w.Code)
	}
}

// TestCreate_InvalidJSON_ContentTypeIsJSON verifies the error response still has JSON Content-Type.
func TestCreate_InvalidJSON_ContentTypeIsJSON(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json on error response, got %q", ct)
	}
}

// TestCreate_InvalidJSON_ReturnsErrorBody verifies the error response body is valid JSON with an error field.
func TestCreate_InvalidJSON_ReturnsErrorBody(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("error response body is not valid JSON: %v", err)
	}
	if resp.Error == "" {
		t.Error("expected non-empty error field in error response")
	}
}

// TestCreate_EmptyBody_Returns400 verifies that an empty request body returns 400.
func TestCreate_EmptyBody_Returns400(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty body, got %d", w.Code)
	}
}

// TestCreate_ValidRequest_ResponseIsValidJSON verifies the success response body is valid JSON.
func TestCreate_ValidRequest_ResponseIsValidJSON(t *testing.T) {
	h := NewTaskHandler(newSuccessStore())
	body := `{"title":"Do laundry","description":"Wash and dry clothes"}`
	w := postJSON(h, body)

	var raw map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&raw); err != nil {
		t.Fatalf("success response body is not valid JSON: %v — body: %q", err, w.Body.String())
	}
}

// TestCreate_ErrorResponses_ContentTypeIsJSON verifies that all 400 error responses set Content-Type.
func TestCreate_ErrorResponses_ContentTypeIsJSON(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"MissingTitle", `{"description":"some desc"}`},
		{"EmptyTitle", `{"title":"","description":"some desc"}`},
		{"LongTitle", fmt.Sprintf(`{"title":%q,"description":"desc"}`, strings.Repeat("x", 201))},
		{"EmptyDescription", `{"title":"Valid","description":""}`},
		{"LongDescription", fmt.Sprintf(`{"title":"Valid","description":%q}`, strings.Repeat("x", 1001))},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewTaskHandler(newSuccessStore())
			w := postJSON(h, tc.body)
			ct := w.Header().Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("expected Content-Type application/json, got %q", ct)
			}
		})
	}
}
