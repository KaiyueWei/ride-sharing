package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)


func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// writeJSONError writes a JSON-encoded error so clients can always parse the
// response body as JSON, regardless of status code.
func writeJSONError(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, contracts.APIResponse{
		Error: &contracts.APIError{
			Code:    http.StatusText(status),
			Message: message,
		},
	})
}