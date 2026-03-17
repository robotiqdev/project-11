package models

import (
	"errors"
	"time"
)

// Task represents the domain task entity.
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest is the DTO for creating a task.
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Validate validates the CreateTaskRequest.
func (r CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Title) > 200 {
		return errors.New("title must not exceed 200 characters")
	}
	if r.Description == "" {
		return errors.New("description is required")
	}
	if len(r.Description) > 1000 {
		return errors.New("description must not exceed 1000 characters")
	}
	return nil
}

// UpdateTaskRequest is the DTO for updating a task.
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Validate validates the UpdateTaskRequest.
func (r UpdateTaskRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Title) > 200 {
		return errors.New("title must not exceed 200 characters")
	}
	if r.Description == "" {
		return errors.New("description is required")
	}
	if len(r.Description) > 1000 {
		return errors.New("description must not exceed 1000 characters")
	}
	return nil
}
