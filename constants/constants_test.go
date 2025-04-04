// File: assignment-2/constants/constants_test.go
// This test file checks that the constants in "constants.go" match the expected values.

package constants

import (
	"testing"
)

// TestConstants verifies that all constants remain at their expected values.

func TestConstants(t *testing.T) {
	// Testing VERSION
	if VERSION != "v1" {
		t.Errorf("Expected VERSION to be 'v1', got '%s'", VERSION)
	}

	// Testing BASE_PATH
	expectedBasePath := "/dashboard/v1/"
	if BASE_PATH != expectedBasePath {
		t.Errorf("Expected BASE_PATH to be '%s', got '%s'", expectedBasePath, BASE_PATH)
	}

	// Testing REGISTRATIONS_PATH
	expectedRegPath := "/dashboard/v1/registrations/"
	if REGISTRATIONS_PATH != expectedRegPath {
		t.Errorf("Expected REGISTRATIONS_PATH to be '%s', got '%s'", expectedRegPath, REGISTRATIONS_PATH)
	}

	// Testing DASHBOARDS_PATH
	expectedDashboardsPath := "/dashboard/v1/dashboards/"
	if DASHBOARDS_PATH != expectedDashboardsPath {
		t.Errorf("Expected DASHBOARDS_PATH to be '%s', got '%s'", expectedDashboardsPath, DASHBOARDS_PATH)
	}

	// Testing NOTIFICATIONS_PATH
	expectedNotificationsPath := "/dashboard/v1/notifications/"
	if NOTIFICATIONS_PATH != expectedNotificationsPath {
		t.Errorf("Expected NOTIFICATIONS_PATH to be '%s', got '%s'", expectedNotificationsPath, NOTIFICATIONS_PATH)
	}

	// Testing STATUS_PATH
	expectedStatusPath := "/dashboard/v1/status/"
	if STATUS_PATH != expectedStatusPath {
		t.Errorf("Expected STATUS_PATH to be '%s', got '%s'", expectedStatusPath, STATUS_PATH)
	}

	// Testing DefaultPort
	if DefaultPort != "8080" {
		t.Errorf("Expected DefaultPort to be '8080', got '%s'", DefaultPort)
	}

	// Testing REST_COUNTRIES_ALPHA
	expectedCountriesAlpha := "http://129.241.150.113:8080/v3.1/alpha/"
	if REST_COUNTRIES_ALPHA != expectedCountriesAlpha {
		t.Errorf("Expected REST_COUNTRIES_ALPHA to be '%s', got '%s'", expectedCountriesAlpha, REST_COUNTRIES_ALPHA)
	}

	// Testing REST_COUNTRIES_NAME
	expectedCountriesName := "http://129.241.150.113:8080/v3.1/name/"
	if REST_COUNTRIES_NAME != expectedCountriesName {
		t.Errorf("Expected REST_COUNTRIES_NAME to be '%s', got '%s'", expectedCountriesName, REST_COUNTRIES_NAME)
	}

	// Testing CURRENCY_API
	expectedCurrencyAPI := "http://129.241.150.113:9090/currency/"
	if CURRENCY_API != expectedCurrencyAPI {
		t.Errorf("Expected CURRENCY_API to be '%s', got '%s'", expectedCurrencyAPI, CURRENCY_API)
	}

	// Testing OPEN_METEO_API
	expectedOpenMeteoAPI := "https://api.open-meteo.com/v1/forecast"
	if OPEN_METEO_API != expectedOpenMeteoAPI {
		t.Errorf("Expected OPEN_METEO_API to be '%s', got '%s'", expectedOpenMeteoAPI, OPEN_METEO_API)
	}

	// Testing Firebase collection names
	if REGISTRATIONS_COLLECTION != "registrations" {
		t.Errorf("Expected REGISTRATIONS_COLLECTION to be 'registrations', got '%s'", REGISTRATIONS_COLLECTION)
	}
	if NOTIFICATIONS_COLLECTION != "notifications" {
		t.Errorf("Expected NOTIFICATIONS_COLLECTION to be 'notifications', got '%s'", NOTIFICATIONS_COLLECTION)
	}
	if CACHE_COLLECTION != "cache" {
		t.Errorf("Expected CACHE_COLLECTION to be 'cache', got '%s'", CACHE_COLLECTION)
	}

	// Testing ServiceVersion
	if ServiceVersion != "v1.0.0" {
		t.Errorf("Expected ServiceVersion to be 'v1.0.0', got '%s'", ServiceVersion)
	}
}
