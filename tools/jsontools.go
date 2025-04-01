// File: assignment-2/tools/jsontools.go
package tools

import (
	"encoding/json"
	"net/http"
)

// WriteJsonResponse marshals 'data' into JSON, sets the response Content-Type,
// writes the status code, and then sends the JSON-encoded data in the body.

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data) // Encode the data to JSON and write it to the response.
	}
}

// WriteJsonErrorResponse writes an error message as JSON using the provided status code.
// It creates a simple JSON object with the key "error" and the given message.

func WriteJsonErrorResponse(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := map[string]string{"error": errMsg}
	json.NewEncoder(w).Encode(resp) // Encode this map to JSON and include it in the response body.
}
