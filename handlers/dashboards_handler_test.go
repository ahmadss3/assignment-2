// File: assignment-2/handlers/dashboards_handler_test.go
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/services"
	"assignment-2/structs"
)

// Backup original references
var (
	origGetRegistrationByID = firebase.GetRegistrationByID
	origFetchCountryInfo    = services.FetchCountryInfo
	origFetchMeteoData      = services.FetchMeteoData
	origFetchCurrencyRates  = services.FetchCurrencyRates
	origTriggerWebhook      = TriggerWebhookEventVar
)

// In-memory store for registrations
var (
	regStore map[string]structs.Registration
	regMutex sync.Mutex
)

// overrideStubs sets up in-memory stubs for firebase + services
func overrideStubs() {
	regStore = make(map[string]structs.Registration)

	// Stub for Firestore: GetRegistrationByID
	firebase.GetRegistrationByID = func(ctx context.Context, docID string) (*structs.Registration, error) {
		regMutex.Lock()
		defer regMutex.Unlock()
		r, ok := regStore[docID]
		if !ok {
			return nil, os.ErrNotExist
		}
		return &r, nil
	}

	// Stub for country info
	services.FetchCountryInfo = func(countryOrISO string) (*structs.CountryInfo, error) {
		if strings.ToUpper(countryOrISO) == "NO" {
			return &structs.CountryInfo{
				Name:         "Norway",
				Capital:      "Oslo",
				Population:   5372000,
				Area:         385207.0,
				BaseCurrency: "NOK",
				Coordinates:  structs.Coordinates{Lat: 60.0, Lon: 10.0},
			}, nil
		} else if strings.ToUpper(countryOrISO) == "ERR" {
			// simulate an error
			return nil, errors.New("failed to fetch country info")
		}
		// fallback
		return &structs.CountryInfo{
			Name:         "SomeCountry",
			Capital:      "SomeCapital",
			Population:   12345,
			Area:         999.9,
			BaseCurrency: "XYZ",
			Coordinates:  structs.Coordinates{Lat: 1.0, Lon: 2.0},
		}, nil
	}

	// Stub for meteo data
	services.FetchMeteoData = func(lat, lon float64) (*structs.MeteoData, error) {
		return &structs.MeteoData{
			AverageTemp:          5.5,
			AveragePrecipitation: 1.2,
		}, nil
	}

	// Stub for currency rates
	services.FetchCurrencyRates = func(base string) (structs.CurrencyRates, error) {
		if base == "NOK" {
			return structs.CurrencyRates{"EUR": 0.09, "USD": 0.1}, nil
		}
		// fallback
		return structs.CurrencyRates{"ABC": 0.5}, nil
	}

	// Stub for TriggerWebhook
	TriggerWebhookEventVar = func(event, country string) {
		// do nothing
	}
}

// revertStubs reverts all stubs to original references
func revertStubs() {
	firebase.GetRegistrationByID = origGetRegistrationByID
	services.FetchCountryInfo = origFetchCountryInfo
	services.FetchMeteoData = origFetchMeteoData
	services.FetchCurrencyRates = origFetchCurrencyRates
	TriggerWebhookEventVar = origTriggerWebhook
}

// Helper to store a registration in the in-memory store
func storeRegistration(docID string, reg structs.Registration) {
	regMutex.Lock()
	defer regMutex.Unlock()
	regStore[docID] = reg
}

// TestDashboardsHandler tests the GET /dashboard/v1/dashboards/{id} route.
func TestDashboardsHandler(t *testing.T) {
	overrideStubs()
	defer revertStubs()

	// Insert one registration with all features
	storeRegistration("doc-1", structs.Registration{
		ID:      "doc-1",
		Country: "Norway",
		ISOCode: "NO",
		Features: structs.Features{
			Temperature:      true,
			Precipitation:    true,
			Capital:          true,
			Coordinates:      true,
			Population:       true,
			Area:             true,
			TargetCurrencies: []string{"EUR", "USD"},
		},
		LastChange: time.Now(),
	})

	t.Run("GetDashboard_Success", func(t *testing.T) {
		// Test a valid doc ID
		req := httptest.NewRequest(http.MethodGet, constants.DASHBOARDS_PATH+"doc-1", nil)
		rr := httptest.NewRecorder()

		DashboardsRouter(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK, got %d", rr.Code)
		}

		var dash structs.Dashboard
		if err := json.Unmarshal(rr.Body.Bytes(), &dash); err != nil {
			t.Fatalf("Failed to parse dashboard JSON: %v", err)
		}

		if dash.Country != "Norway" {
			t.Errorf("Expected Country=Norway, got %s", dash.Country)
		}
		if dash.Features.Temperature != 5.5 {
			t.Errorf("Expected Temperature=5.5, got %f", dash.Features.Temperature)
		}
		if dash.Features.TargetCurrencies["EUR"] != 0.09 {
			t.Errorf("Expected EUR=0.09, got %f", dash.Features.TargetCurrencies["EUR"])
		}
	})

	t.Run("GetDashboard_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, constants.DASHBOARDS_PATH+"doc-999", nil)
		rr := httptest.NewRecorder()

		DashboardsRouter(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected 404 Not Found, got %d", rr.Code)
		}
	})

	t.Run("GetDashboard_CountryInfoError", func(t *testing.T) {
		// Insert a registration referencing ISOCode=ERR => triggers error in stub
		storeRegistration("doc-err", structs.Registration{
			ID:      "doc-err",
			Country: "",
			ISOCode: "ERR",
			Features: structs.Features{
				Capital:    true,
				Population: true,
			},
		})

		req := httptest.NewRequest(http.MethodGet, constants.DASHBOARDS_PATH+"doc-err", nil)
		rr := httptest.NewRecorder()
		DashboardsRouter(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected 200 OK (with partial data), got %d", rr.Code)
		}
	})

	t.Run("MethodNotAllowed_Post", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, constants.DASHBOARDS_PATH+"doc-1", nil)
		rr := httptest.NewRecorder()

		DashboardsRouter(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected 405 Method Not Allowed, got %d", rr.Code)
		}
	})

	t.Run("MethodNotAllowed_ListAll", func(t *testing.T) {
		// If path == /dashboard/v1/dashboards/ with no ID
		req := httptest.NewRequest(http.MethodGet, constants.DASHBOARDS_PATH, nil)
		rr := httptest.NewRecorder()

		DashboardsRouter(rr, req)
		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected 405 because listing is not allowed, got %d", rr.Code)
		}
	})
}
