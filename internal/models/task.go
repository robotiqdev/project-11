package models

import "time"

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

// Validate checks that the request fields satisfy all constraints.
func (r *CreateTaskRequest) Validate() error {
	return nil // stub — implementation pending
}
