package router

import (
	"net/http"

	"github.com/workspace/repo/internal/storage"
)

// New constructs an http.Handler with all task routes registered.
func New(store storage.TaskStore) http.Handler {
	// stub — implementation pending
	_ = store
	return http.NewServeMux()
}
