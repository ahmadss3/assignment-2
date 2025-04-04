// File: assignment-2/services/status_check_test.go
package services

import (
	"errors"
	"net/http"
	"testing"
)

// TestCheckCountriesAPI verifies the integration call to the REST Countries endpoint.
func TestCheckCountriesAPI(t *testing.T) {
	result := CheckCountriesAPI()
	if result.Error != nil {
		t.Logf("CheckCountriesAPI encountered an error: %v", result.Error)
		t.Fail()
	} else {
		// Check if status code is in a valid range.
		if result.StatusCode < 200 || result.StatusCode > 399 {
			t.Errorf("CheckCountriesAPI returned status code %d, expected a 2xx or 3xx range", result.StatusCode)
		}
	}
}

// TestCheckOpenMeteo verifies the integration call to Open-Meteo with a minimal query.
func TestCheckOpenMeteo(t *testing.T) {
	result := CheckOpenMeteo()
	if result.Error != nil {
		t.Logf("CheckOpenMeteo encountered an error: %v", result.Error)
		t.Fail()
	} else {
		if result.StatusCode < 200 || result.StatusCode > 399 {
			t.Errorf("CheckOpenMeteo returned status code %d, expected 2xx or 3xx", result.StatusCode)
		}
	}
}

// TestCheckCurrencyAPI verifies the integration call to the currency API.
func TestCheckCurrencyAPI(t *testing.T) {
	result := CheckCurrencyAPI()
	if result.Error != nil {
		t.Logf("CheckCurrencyAPI encountered an error: %v", result.Error)
		t.Fail()
	} else {
		if result.StatusCode < 200 || result.StatusCode > 399 {
			t.Errorf("CheckCurrencyAPI returned status code %d, expected 2xx or 3xx", result.StatusCode)
		}
	}
}

// TestTranslateErrorToStatus checks that errors map to 503, and nil maps to 200.
func TestTranslateErrorToStatus(t *testing.T) {
	if code := TranslateErrorToStatus(nil); code != http.StatusOK {
		t.Errorf("Expected TranslateErrorToStatus(nil) = 200, got %d", code)
	}
	dummyErr := errors.New("some network error")
	if code := TranslateErrorToStatus(dummyErr); code != http.StatusServiceUnavailable {
		t.Errorf("Expected TranslateErrorToStatus(dummyErr) = 503, got %d", code)
	}
}

// TestCheckFirebaseNotifications verifies the placeholder functionality that returns 200 if no error from GetAllNotifications.
func TestCheckFirebaseNotifications(t *testing.T) {
	result := CheckFirebaseNotifications()
	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}
}
