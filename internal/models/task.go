package models

import "time"

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
	panic("not implemented")
}

// UpdateTaskRequest is the DTO for updating a task.
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Validate validates the UpdateTaskRequest.
func (r UpdateTaskRequest) Validate() error {
	panic("not implemented")
}
