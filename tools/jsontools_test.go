// File: assignment-2/tools/jsontools_test.go
package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// TestWriteJsonResponse checks whether WriteJsonResponse correctly sets
// the response headers, status code, and JSON body based on the input data.

func TestWriteJsonResponse(t *testing.T) {
	t.Run("NonNilData", func(t *testing.T) {
		// Create a ResponseRecorder to record the response.
		rr := httptest.NewRecorder()

		// Prepare some data to be serialized to JSON.
		sampleData := map[string]interface{}{
			"message": "Test successful",
			"value":   42,
		}
		// Call WriteJsonResponse with status code 200 and the sample data.
		WriteJsonResponse(rr, http.StatusOK, sampleData)

		// Verify status code.
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		// Verify Content-Type.
		if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", ct)
		}

		// Parse the response body as JSON.
		var parsed map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &parsed); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		expected := map[string]interface{}{
			"message": "Test successful",
			"value":   float64(42), // JSON unmarshals numbers to float64 by default
		}
		if !reflect.DeepEqual(expected, parsed) {
			t.Errorf("Body mismatch.\nExpected: %v\nGot:      %v", expected, parsed)
		}
	})

	t.Run("NilData", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// Call WriteJsonResponse with nil data.
		WriteJsonResponse(rr, http.StatusNoContent, nil)

		// Verify status code.
		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
		}

		// Verify Content-Type.
		if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", ct)
		}

		// The body should be empty if data is nil.
		if body := rr.Body.String(); strings.TrimSpace(body) != "" {
			t.Errorf("Expected empty body, got '%s'", body)
		}
	})
}

// TestWriteJsonErrorResponse checks whether WriteJsonErrorResponse correctly sets
// the response headers, status code, and writes a JSON object with "error".
func TestWriteJsonErrorResponse(t *testing.T) {
	t.Run("BasicError", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// Call WriteJsonErrorResponse with a custom status code and message.
		WriteJsonErrorResponse(rr, http.StatusBadRequest, "Invalid request")

		// Verify status code.
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		// Verify Content-Type.
		if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", ct)
		}

		// Parse JSON body.
		var parsed map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &parsed); err != nil {
			t.Fatalf("Failed to unmarshal JSON error response: %v", err)
		}

		// Expect a single key "error" with the custom message.
		expected := map[string]string{"error": "Invalid request"}
		if !reflect.DeepEqual(expected, parsed) {
			t.Errorf("Error response mismatch.\nExpected: %v\nGot:      %v", expected, parsed)
		}
	})

	t.Run("EmptyErrorMsg", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// If we send an empty error message.
		WriteJsonErrorResponse(rr, http.StatusInternalServerError, "")

		// Check status code.
		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}

		// Check content type.
		if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Expected 'application/json', got '%s'", ct)
		}

		// Unmarshal response body.
		var parsed map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &parsed); err != nil {
			t.Fatalf("Failed to unmarshal JSON error response: %v", err)
		}

		// Expect "error" key, but empty message.
		expected := map[string]string{"error": ""}
		if !reflect.DeepEqual(expected, parsed) {
			t.Errorf("Error response mismatch.\nExpected: %v\nGot:      %v", expected, parsed)
		}
	})
}
