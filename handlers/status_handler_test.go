// File: assignment-2/handlers/status_handler_test.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
)

// TestMain runs once for this package, letting us init Firebase for integration tests.
func TestMain(m *testing.M) {
	err := firebase.InitFirebase()
	if err != nil {
		log.Printf("Warning: could not init Firebase in test: %v\n", err)
		os.Exit(0)
	}

	code := m.Run()

	// Close Firestore
	if firebase.FirestoreClient != nil {
		_ = firebase.FirestoreClient.Close()
	}

	os.Exit(code)
}

// TestStatusHandler_MethodNotAllowed ensures a POST returns 405.
func TestStatusHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, constants.STATUS_PATH, nil)
	rr := httptest.NewRecorder()

	StatusHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse JSON error: %v", err)
	} else if resp["error"] == "" {
		t.Error("Expected an 'error' field in JSON response")
	}
}

// TestStatusHandler_Get checks a GET request on /status.
func TestStatusHandler_Get(t *testing.T) {
	// Use a start time slightly in the past so uptime is positive
	AssignStartTime(time.Now().Add(-5 * time.Second))

	req := httptest.NewRequest(http.MethodGet, constants.STATUS_PATH, nil)
	rr := httptest.NewRecorder()

	StatusHandler(rr, req)

	// Typically 200 if external calls + DB are OK, 503 if something fails
	if rr.Code != http.StatusOK && rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 200 or 503, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Could not parse JSON: %v", err)
	}

	fields := []string{
		"countries_api",
		"meteo_api",
		"currency_api",
		"notification_db",
		"webhooks",
		"version",
		"uptime",
	}
	for _, f := range fields {
		if _, ok := resp[f]; !ok {
			t.Errorf("Expected '%s' in status JSON, missing", f)
		}
	}

	// Check uptime >= 0
	if val, ok := resp["uptime"].(float64); ok {
		if val < 0 {
			t.Errorf("Expected uptime >= 0, got %f", val)
		}
	} else {
		t.Error("Expected 'uptime' to be a float64")
	}
}

// TestStatusHandler_JSONContents checks more specific fields in the JSON.
func TestStatusHandler_JSONContents(t *testing.T) {
	AssignStartTime(time.Now().Add(-10 * time.Second))

	req := httptest.NewRequest(http.MethodGet, constants.STATUS_PATH, nil)
	rr := httptest.NewRecorder()

	StatusHandler(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 200 or 503, got %d", rr.Code)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// countries_api => 200 or 503
	if val, ok := parsed["countries_api"].(float64); ok {
		if val != 200 && val != 503 {
			t.Errorf("countries_api is %v, expected 200 or 503", val)
		}
	} else {
		t.Error("Missing or invalid type for 'countries_api'")
	}

	// notification_db => 200 or 503
	if val, ok := parsed["notification_db"].(float64); ok {
		if val != 200 && val != 503 {
			t.Errorf("notification_db is %v, expected 200 or 503", val)
		}
	} else {
		t.Error("Missing or invalid type for 'notification_db'")
	}
}
