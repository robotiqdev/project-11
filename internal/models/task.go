package models

import (
	"errors"
	"time"
)

// Task represents a task in the system.
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTaskRequest holds the fields required to create a new task.
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTaskRequest holds the fields that can be updated on an existing task.
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Validate checks that the request fields satisfy all constraints.
func (r *CreateTaskRequest) Validate() error {
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
