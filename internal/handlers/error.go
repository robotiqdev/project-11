package handlers

import "net/http"

// ErrorResponse is the JSON shape returned for error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeError is a stub — implementation pending.
func writeError(w http.ResponseWriter, status int, message string) {
	// not yet implemented
}
