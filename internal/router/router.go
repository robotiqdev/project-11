package router

import (
	"net/http"

	"github.com/workspace/repo/internal/handlers"
	"github.com/workspace/repo/internal/storage"
)

// New constructs an http.Handler with all task routes registered.
func New(store storage.TaskStore) http.Handler {
	mux := http.NewServeMux()
	handler := handlers.NewTaskHandler(store)

	mux.HandleFunc("GET /tasks", handler.GetAll)
	mux.HandleFunc("POST /tasks", handler.Create)
	mux.HandleFunc("GET /tasks/{id}", handler.GetByID)
	mux.HandleFunc("PUT /tasks/{id}", handler.Update)
	mux.HandleFunc("DELETE /tasks/{id}", handler.Delete)

	return mux
}
