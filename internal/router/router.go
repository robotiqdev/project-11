package router

import (
	"net/http"

	"github.com/workspace/repo/internal/handlers"
)

// New creates and returns a new HTTP ServeMux with all task routes registered.
func New(h *handlers.TaskHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", h.GetByID)
	mux.HandleFunc("DELETE /tasks/{id}", h.Delete)
	return mux
}
