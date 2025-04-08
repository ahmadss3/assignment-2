// file assignment-2/handlers/registrations_handler_test.go

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/structs"
)

// In-memory "registrations" collection:
// docID -> structs.Registration
var (
	stubRegMutex sync.Mutex
	stubRegStore = make(map[string]structs.Registration)
	idCounter    int
)

// Keep a backup of the original firebase.*Registration function variables.
var (
	originalSaveRegistration    = firebase.SaveRegistration
	originalGetAllRegistrations = firebase.GetAllRegistrations
	originalGetRegistrationByID = firebase.GetRegistrationByID
	originalUpdateRegistration  = firebase.UpdateRegistration
	originalDeleteRegistration  = firebase.DeleteRegistration
	originalPatchRegistration   = firebase.PatchRegistration
)

// overrideFirebaseStubs redirects the firebase.*Registration variables to in-memory stub implementations.
func overrideFirebaseStubs() {
	firebase.SaveRegistration = func(ctx context.Context, reg structs.Registration) (string, error) {
		stubRegMutex.Lock()
		defer stubRegMutex.Unlock()
		idCounter++
		docID := "doc-" + string(rune(idCounter))
		reg.ID = docID
		stubRegStore[docID] = reg
		return docID, nil
	}

	firebase.GetAllRegistrations = func(ctx context.Context) ([]structs.Registration, error) {
		stubRegMutex.Lock()
		defer stubRegMutex.Unlock()
		var all []structs.Registration
		for _, r := range stubRegStore {
			all = append(all, r)
		}
		return all, nil
	}

	firebase.GetRegistrationByID = func(ctx context.Context, docID string) (*structs.Registration, error) {
		stubRegMutex.Lock()
		defer stubRegMutex.Unlock()
		reg, ok := stubRegStore[docID]
		if !ok {
			return nil, os.ErrNotExist
		}
		return &reg, nil
	}

	firebase.UpdateRegistration = func(ctx context.Context, docID string, reg structs.Registration) error {
		stubRegMutex.Lock()
		defer stubRegMutex.Unlock()
		_, exists := stubRegStore[docID]
		if !exists {
			return os.ErrNotExist
		}
		reg.ID = docID
		stubRegStore[docID] = reg
		return nil
	}

	firebase.DeleteRegistration = func(ctx context.Context, docID string) error {
		stubRegMutex.Lock()
		defer stubRegMutex.Unlock()
		_, exists := stubRegStore[docID]
		if !exists {
			return os.ErrNotExist
		}
		delete(stubRegStore, docID)
		return nil
	}

	firebase.PatchRegistration = func(ctx context.Context, docID string, partial structs.Registration) error {
		stubRegMutex.Lock()
		defer stubRegMutex.Unlock()
		existing, ok := stubRegStore[docID]
		if !ok {
			return os.ErrNotExist
		}
		// Minimal patch-like logic:
		if partial.Country != "" {
			existing.Country = partial.Country
		}
		if partial.ISOCode != "" {
			existing.ISOCode = partial.ISOCode
		}

		if !isZeroFeatures(partial.Features) {
			existing.Features = partial.Features
		}
		existing.LastChange = time.Now()
		existing.ID = docID
		stubRegStore[docID] = existing
		return nil
	}
}

// revertFirebaseStubs restores the original references
func revertFirebaseStubs() {
	firebase.SaveRegistration = originalSaveRegistration
	firebase.GetAllRegistrations = originalGetAllRegistrations
	firebase.GetRegistrationByID = originalGetRegistrationByID
	firebase.UpdateRegistration = originalUpdateRegistration
	firebase.DeleteRegistration = originalDeleteRegistration
	firebase.PatchRegistration = originalPatchRegistration
}

// isZeroFeatures checks if a Features struct is entirely "zero".

func isZeroFeatures(f structs.Features) bool {
	return !f.Temperature &&
		!f.Precipitation &&
		!f.Capital &&
		!f.Coordinates &&
		!f.Population &&
		!f.Area &&
		len(f.TargetCurrencies) == 0
}

// TestRegistrationsHandler runs subtests for POST, GET, PUT, PATCH, DELETE
func TestRegistrationsHandler(t *testing.T) {
	// Override stubs at the start of this test function
	overrideFirebaseStubs()
	defer revertFirebaseStubs()

	t.Run("PostRegistration_Success", func(t *testing.T) {
		body := `{
          "country": "TestCountry",
          "isoCode": "TC",
          "features": {
            "temperature": true,
            "precipitation": false,
            "capital": true,
            "coordinates": false,
            "population": false,
            "area": false,
            "targetCurrencies": ["USD","EUR"]
          }
        }`

		req := httptest.NewRequest(http.MethodPost, constants.REGISTRATIONS_PATH, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		RegistrationRouter(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected 201, got %d", rr.Code)
		}
		var resp map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}
		if resp["id"] == nil {
			t.Error("Expected an 'id' field in response")
		}
	})

	t.Run("PostRegistration_InvalidJSON", func(t *testing.T) {
		body := `{"country":"Invalid`
		req := httptest.NewRequest(http.MethodPost, constants.REGISTRATIONS_PATH, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		RegistrationRouter(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", rr.Code)
		}
	})

	t.Run("GetAllRegistrations", func(t *testing.T) {
		// Create a doc
		docID := createFakeRegistration(t, "AllRegTest")
		_ = docID

		req := httptest.NewRequest(http.MethodGet, constants.REGISTRATIONS_PATH, nil)
		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK, got %d", rr.Code)
		}
		var list []map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
			t.Errorf("Failed to parse list JSON: %v", err)
		}
	})

	t.Run("GetOne_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, constants.REGISTRATIONS_PATH+"doc-9999", nil)
		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404 Not Found, got %d", rr.Code)
		}
	})

	t.Run("PutRegistration_Success", func(t *testing.T) {
		docID := createFakeRegistration(t, "PutCountry")
		putBody := `{
          "country": "UpdatedCountry",
          "isoCode": "UC",
          "features": {
            "temperature": false,
            "precipitation": true
          }
        }`
		req := httptest.NewRequest(http.MethodPut, constants.REGISTRATIONS_PATH+docID, strings.NewReader(putBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected 204 No Content, got %d", rr.Code)
		}
	})

	t.Run("PutRegistration_NotFound", func(t *testing.T) {
		putBody := `{
          "country": "NoDoc",
          "isoCode": "ND",
          "features": { "temperature": true }
        }`
		req := httptest.NewRequest(http.MethodPut, constants.REGISTRATIONS_PATH+"doc-9999", strings.NewReader(putBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404 Not Found, got %d", rr.Code)
		}
	})

	t.Run("PatchRegistration_Success", func(t *testing.T) {
		docID := createFakeRegistration(t, "PatchMe")
		patchBody := `{
          "isoCode":"PATCHED",
          "features":{"capital":true}
        }`
		req := httptest.NewRequest(http.MethodPatch, constants.REGISTRATIONS_PATH+docID, strings.NewReader(patchBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected 204, got %d", rr.Code)
		}
	})

	t.Run("PatchRegistration_InvalidJSON", func(t *testing.T) {
		docID := createFakeRegistration(t, "PatchFail")
		body := `{"features": { "temperature":`
		req := httptest.NewRequest(http.MethodPatch, constants.REGISTRATIONS_PATH+docID, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
		}
	})

	t.Run("PatchRegistration_NotFound", func(t *testing.T) {
		body := `{"isoCode":"ABC"}`
		req := httptest.NewRequest(http.MethodPatch, constants.REGISTRATIONS_PATH+"doc-9999", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404 Not Found, got %d", rr.Code)
		}
	})

	t.Run("DeleteRegistration_Success", func(t *testing.T) {
		docID := createFakeRegistration(t, "DelMe")
		req := httptest.NewRequest(http.MethodDelete, constants.REGISTRATIONS_PATH+docID, nil)
		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected 204 No Content, got %d", rr.Code)
		}
	})

	t.Run("DeleteRegistration_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, constants.REGISTRATIONS_PATH+"doc-9999", nil)
		rr := httptest.NewRecorder()
		RegistrationRouter(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404 Not Found for unknown doc, got %d", rr.Code)
		}
	})
}

// createFakeRegistration is a helper that does a POST /registrations/ to create a doc
// in in-memory stub storage, returning the docID from the response.
func createFakeRegistration(t *testing.T, countryName string) string {
	body := `{
      "country": "` + countryName + `",
      "isoCode": "FAKE",
      "features": {
        "temperature": true,
        "precipitation": false
      }
    }`
	req := httptest.NewRequest(http.MethodPost, constants.REGISTRATIONS_PATH, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	RegistrationRouter(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201, got %d", rr.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	docID, ok := resp["id"].(string)
	if !ok {
		t.Fatalf("No 'id' field or not a string in response: %v", resp["id"])
	}
	return docID
}
