package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the JSON shape returned for error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeError writes a JSON error response with the given status code and message.
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
