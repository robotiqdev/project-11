package storage

import "errors"

// ErrNotFound is returned by store operations when a task does not exist.
var ErrNotFound = errors.New("task not found")
