package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// --- CreateTaskRequest.Validate() tests ---

func TestCreateTaskRequest_Validate_EmptyTitleReturnsError(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "",
		Description: "A valid description",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestCreateTaskRequest_Validate_TitleOver200CharsReturnsError(t *testing.T) {
	req := CreateTaskRequest{
		Title:       strings.Repeat("a", 201),
		Description: "A valid description",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for title over 200 chars, got nil")
	}
}

func TestCreateTaskRequest_Validate_TitleExactly200CharsIsValid(t *testing.T) {
	req := CreateTaskRequest{
		Title:       strings.Repeat("a", 200),
		Description: "A valid description",
	}
	err := req.Validate()
	if err != nil {
		t.Fatalf("expected no error for title of exactly 200 chars, got: %v", err)
	}
}

func TestCreateTaskRequest_Validate_EmptyDescriptionReturnsError(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "Valid Title",
		Description: "",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty description, got nil")
	}
}

func TestCreateTaskRequest_Validate_DescriptionOver1000CharsReturnsError(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "Valid Title",
		Description: strings.Repeat("d", 1001),
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for description over 1000 chars, got nil")
	}
}

func TestCreateTaskRequest_Validate_DescriptionExactly1000CharsIsValid(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "Valid Title",
		Description: strings.Repeat("d", 1000),
	}
	err := req.Validate()
	if err != nil {
		t.Fatalf("expected no error for description of exactly 1000 chars, got: %v", err)
	}
}

func TestCreateTaskRequest_Validate_ValidRequestPassesValidation(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "A valid title",
		Description: "A valid description",
	}
	err := req.Validate()
	if err != nil {
		t.Fatalf("expected no error for valid request, got: %v", err)
	}
}

func TestCreateTaskRequest_Validate_ZeroValueIsInvalid(t *testing.T) {
	var req CreateTaskRequest
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for zero-value CreateTaskRequest, got nil")
	}
}

func TestCreateTaskRequest_Validate_ErrorMessageDescribesMissingTitle(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "",
		Description: "A valid description",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
	if err.Error() == "" {
		t.Fatal("expected a non-empty error message for missing title")
	}
}

func TestCreateTaskRequest_Validate_ErrorMessageDescribesTitleTooLong(t *testing.T) {
	req := CreateTaskRequest{
		Title:       strings.Repeat("x", 201),
		Description: "A valid description",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for title over 200 chars, got nil")
	}
	if err.Error() == "" {
		t.Fatal("expected a non-empty error message for title too long")
	}
}

func TestCreateTaskRequest_Validate_ErrorMessageDescribesMissingDescription(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "Valid Title",
		Description: "",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty description, got nil")
	}
	if err.Error() == "" {
		t.Fatal("expected a non-empty error message for missing description")
	}
}

func TestCreateTaskRequest_Validate_ErrorMessagedescribesDescriptionTooLong(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "Valid Title",
		Description: strings.Repeat("d", 1001),
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for description over 1000 chars, got nil")
	}
	if err.Error() == "" {
		t.Fatal("expected a non-empty error message for description too long")
	}
}

// --- UpdateTaskRequest.Validate() tests ---

func TestUpdateTaskRequest_Validate_EmptyTitleReturnsError(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       "",
		Description: "A valid description",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestUpdateTaskRequest_Validate_TitleOver200CharsReturnsError(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       strings.Repeat("a", 201),
		Description: "A valid description",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for title over 200 chars, got nil")
	}
}

func TestUpdateTaskRequest_Validate_TitleExactly200CharsIsValid(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       strings.Repeat("a", 200),
		Description: "A valid description",
	}
	err := req.Validate()
	if err != nil {
		t.Fatalf("expected no error for title of exactly 200 chars, got: %v", err)
	}
}

func TestUpdateTaskRequest_Validate_EmptyDescriptionReturnsError(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       "Valid Title",
		Description: "",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty description, got nil")
	}
}

func TestUpdateTaskRequest_Validate_DescriptionOver1000CharsReturnsError(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       "Valid Title",
		Description: strings.Repeat("d", 1001),
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for description over 1000 chars, got nil")
	}
}

func TestUpdateTaskRequest_Validate_DescriptionExactly1000CharsIsValid(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       "Valid Title",
		Description: strings.Repeat("d", 1000),
	}
	err := req.Validate()
	if err != nil {
		t.Fatalf("expected no error for description of exactly 1000 chars, got: %v", err)
	}
}

func TestUpdateTaskRequest_Validate_ValidRequestPassesValidation(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       "A valid title",
		Description: "A valid description",
	}
	err := req.Validate()
	if err != nil {
		t.Fatalf("expected no error for valid request, got: %v", err)
	}
}

func TestUpdateTaskRequest_Validate_ZeroValueIsInvalid(t *testing.T) {
	var req UpdateTaskRequest
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for zero-value UpdateTaskRequest, got nil")
	}
}

func TestUpdateTaskRequest_Validate_TitleExactly201CharsReturnsError(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       strings.Repeat("a", 201),
		Description: "A valid description",
	}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for title of exactly 201 chars, got nil")
	}
}

func TestUpdateTaskRequest_Validate_ErrorMessageForEmptyTitleMatchesCreate(t *testing.T) {
	createReq := CreateTaskRequest{Title: "", Description: "A valid description"}
	updateReq := UpdateTaskRequest{Title: "", Description: "A valid description"}

	createErr := createReq.Validate()
	updateErr := updateReq.Validate()

	if createErr == nil {
		t.Fatal("expected CreateTaskRequest to return error for empty title")
	}
	if updateErr == nil {
		t.Fatal("expected UpdateTaskRequest to return error for empty title")
	}
	if createErr.Error() != updateErr.Error() {
		t.Errorf("error messages differ: CreateTaskRequest=%q, UpdateTaskRequest=%q",
			createErr.Error(), updateErr.Error())
	}
}

func TestUpdateTaskRequest_Validate_ErrorMessageForTitleTooLongMatchesCreate(t *testing.T) {
	createReq := CreateTaskRequest{Title: strings.Repeat("a", 201), Description: "A valid description"}
	updateReq := UpdateTaskRequest{Title: strings.Repeat("a", 201), Description: "A valid description"}

	createErr := createReq.Validate()
	updateErr := updateReq.Validate()

	if createErr == nil {
		t.Fatal("expected CreateTaskRequest to return error for title over 200 chars")
	}
	if updateErr == nil {
		t.Fatal("expected UpdateTaskRequest to return error for title over 200 chars")
	}
	if createErr.Error() != updateErr.Error() {
		t.Errorf("error messages differ: CreateTaskRequest=%q, UpdateTaskRequest=%q",
			createErr.Error(), updateErr.Error())
	}
}

func TestUpdateTaskRequest_Validate_ErrorMessageForEmptyDescriptionMatchesCreate(t *testing.T) {
	createReq := CreateTaskRequest{Title: "Valid Title", Description: ""}
	updateReq := UpdateTaskRequest{Title: "Valid Title", Description: ""}

	createErr := createReq.Validate()
	updateErr := updateReq.Validate()

	if createErr == nil {
		t.Fatal("expected CreateTaskRequest to return error for empty description")
	}
	if updateErr == nil {
		t.Fatal("expected UpdateTaskRequest to return error for empty description")
	}
	if createErr.Error() != updateErr.Error() {
		t.Errorf("error messages differ: CreateTaskRequest=%q, UpdateTaskRequest=%q",
			createErr.Error(), updateErr.Error())
	}
}

func TestUpdateTaskRequest_Validate_ErrorMessageForDescriptionTooLongMatchesCreate(t *testing.T) {
	createReq := CreateTaskRequest{Title: "Valid Title", Description: strings.Repeat("d", 1001)}
	updateReq := UpdateTaskRequest{Title: "Valid Title", Description: strings.Repeat("d", 1001)}

	createErr := createReq.Validate()
	updateErr := updateReq.Validate()

	if createErr == nil {
		t.Fatal("expected CreateTaskRequest to return error for description over 1000 chars")
	}
	if updateErr == nil {
		t.Fatal("expected UpdateTaskRequest to return error for description over 1000 chars")
	}
	if createErr.Error() != updateErr.Error() {
		t.Errorf("error messages differ: CreateTaskRequest=%q, UpdateTaskRequest=%q",
			createErr.Error(), updateErr.Error())
	}
}

func TestUpdateTaskRequest_Validate_ExactErrorMessageEmptyTitle(t *testing.T) {
	req := UpdateTaskRequest{Title: "", Description: "A valid description"}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
	want := "title is required"
	if err.Error() != want {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}

func TestUpdateTaskRequest_Validate_ExactErrorMessageTitleTooLong(t *testing.T) {
	req := UpdateTaskRequest{Title: strings.Repeat("a", 201), Description: "A valid description"}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for title over 200 chars, got nil")
	}
	want := "title must not exceed 200 characters"
	if err.Error() != want {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}

func TestUpdateTaskRequest_Validate_ExactErrorMessageEmptyDescription(t *testing.T) {
	req := UpdateTaskRequest{Title: "Valid Title", Description: ""}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for empty description, got nil")
	}
	want := "description is required"
	if err.Error() != want {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}

func TestUpdateTaskRequest_Validate_ExactErrorMessageDescriptionTooLong(t *testing.T) {
	req := UpdateTaskRequest{Title: "Valid Title", Description: strings.Repeat("d", 1001)}
	err := req.Validate()
	if err == nil {
		t.Fatal("expected error for description over 1000 chars, got nil")
	}
	want := "description must not exceed 1000 characters"
	if err.Error() != want {
		t.Errorf("unexpected error message: got %q, want %q", err.Error(), want)
	}
}

// --- Task struct tests ---

func TestTask_ZeroValueIsInvalid(t *testing.T) {
	var task Task
	// A zero-value Task should have empty Title and Description.
	// Validate is defined on request DTOs; zero-value Task has no meaningful title/description.
	if task.Title != "" {
		t.Errorf("expected zero-value Task.Title to be empty, got %q", task.Title)
	}
	if task.Description != "" {
		t.Errorf("expected zero-value Task.Description to be empty, got %q", task.Description)
	}
	if task.ID != 0 {
		t.Errorf("expected zero-value Task.ID to be 0, got %d", task.ID)
	}
	if !task.CreatedAt.IsZero() {
		t.Errorf("expected zero-value Task.CreatedAt to be zero time, got %v", task.CreatedAt)
	}
	if !task.UpdatedAt.IsZero() {
		t.Errorf("expected zero-value Task.UpdatedAt to be zero time, got %v", task.UpdatedAt)
	}
}

// --- JSON serialization round-trip tests ---

func TestTask_JSONRoundTrip(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	original := Task{
		ID:          42,
		Title:       "Test Task",
		Description: "A description for testing",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Task
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %d, want %d", decoded.ID, original.ID)
	}
	if decoded.Title != original.Title {
		t.Errorf("Title mismatch: got %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Description != original.Description {
		t.Errorf("Description mismatch: got %q, want %q", decoded.Description, original.Description)
	}
	if !decoded.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", decoded.CreatedAt, original.CreatedAt)
	}
	if !decoded.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", decoded.UpdatedAt, original.UpdatedAt)
	}
}

func TestTask_JSONUsesSnakeCaseKeys(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	task := Task{
		ID:          1,
		Title:       "Title",
		Description: "Desc",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal to map failed: %v", err)
	}

	expectedKeys := []string{"id", "title", "description", "created_at", "updated_at"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q to be present, but it was not. JSON: %s", key, string(data))
		}
	}
}

func TestCreateTaskRequest_JSONRoundTrip(t *testing.T) {
	original := CreateTaskRequest{
		Title:       "My Task",
		Description: "My Description",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded CreateTaskRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Title != original.Title {
		t.Errorf("Title mismatch: got %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Description != original.Description {
		t.Errorf("Description mismatch: got %q, want %q", decoded.Description, original.Description)
	}
}

func TestUpdateTaskRequest_JSONRoundTrip(t *testing.T) {
	original := UpdateTaskRequest{
		Title:       "Updated Task",
		Description: "Updated Description",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded UpdateTaskRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Title != original.Title {
		t.Errorf("Title mismatch: got %q, want %q", decoded.Title, original.Title)
	}
	if decoded.Description != original.Description {
		t.Errorf("Description mismatch: got %q, want %q", decoded.Description, original.Description)
	}
}

func TestCreateTaskRequest_JSONUsesSnakeCaseKeys(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "Title",
		Description: "Desc",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal to map failed: %v", err)
	}

	for _, key := range []string{"title", "description"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q, not found in: %s", key, string(data))
		}
	}
}

func TestUpdateTaskRequest_JSONUsesSnakeCaseKeys(t *testing.T) {
	req := UpdateTaskRequest{
		Title:       "Title",
		Description: "Desc",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal to map failed: %v", err)
	}

	for _, key := range []string{"title", "description"} {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected JSON key %q, not found in: %s", key, string(data))
		}
	}
}
