package structs

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestExternalStructsDefaults(t *testing.T) {
	t.Run("TestCountryInfoDefaults", func(t *testing.T) {
		var c CountryInfo
		if c.Name != "" || c.Capital != "" || c.Population != 0 || c.Area != 0 || c.BaseCurrency != "" {
			t.Error("Expected zero/default values in CountryInfo")
		}
		if c.Coordinates.Lat != 0 || c.Coordinates.Lon != 0 {
			t.Error("Expected zero lat/lon by default")
		}
	})

	t.Run("TestMeteoDataDefaults", func(t *testing.T) {
		var m MeteoData
		if m.AverageTemp != 0 || m.AveragePrecipitation != 0 {
			t.Error("Expected zero defaults in MeteoData")
		}
	})

	t.Run("TestCurrencyRatesDefaults", func(t *testing.T) {
		var cr CurrencyRates
		if len(cr) != 0 {
			t.Errorf("Expected empty map, got size %d", len(cr))
		}
	})

	t.Run("TestCacheEntryExDefaults", func(t *testing.T) {
		var ce CacheEntryEx
		if ce.Key != "" || ce.Data != nil || !ce.LastFetched.IsZero() || ce.TTLHours != 0 {
			t.Error("Expected zero/default values in CacheEntryEx")
		}
	})
}

func TestCacheEntryExFields(t *testing.T) {
	now := time.Now()
	ce := CacheEntryEx{
		Key:         "testKey",
		Data:        []byte(`{"foo":"bar"}`),
		LastFetched: now,
		TTLHours:    24,
	}
	if ce.Key != "testKey" {
		t.Errorf("Expected Key='testKey', got %s", ce.Key)
	}
	if string(ce.Data) != `{"foo":"bar"}` {
		t.Errorf("Expected Data={\"foo\":\"bar\"}, got %s", ce.Data)
	}
	if ce.TTLHours != 24 {
		t.Errorf("Expected TTLHours=24, got %d", ce.TTLHours)
	}
	if !ce.LastFetched.Equal(now) {
		t.Error("LastFetched mismatch from assignment.")
	}
}

// TestDashboardStruct merges the old "dashboard_test.go"
func TestDashboardStruct(t *testing.T) {
	t.Run("FieldAssignments", func(t *testing.T) {
		dash := Dashboard{
			Country:       "Norway",
			ISOCode:       "NO",
			LastRetrieval: time.Date(2025, 4, 1, 12, 30, 0, 0, time.UTC),
			Features: DashboardFeatures{
				Temperature:   -2.5,
				Precipitation: 1.2,
				Capital:       "Oslo",
				Coordinates: &Coordinates{
					Lat: 59.95, Lon: 10.75,
				},
				Population:       5370000,
				Area:             323802,
				TargetCurrencies: map[string]float64{"EUR": 0.095, "USD": 0.10},
			},
		}
		if dash.Country != "Norway" {
			t.Errorf("Expected Country='Norway', got '%s'", dash.Country)
		}
		if dash.ISOCode != "NO" {
			t.Errorf("Expected ISOCode='NO', got '%s'", dash.ISOCode)
		}
		expectedTime := time.Date(2025, 4, 1, 12, 30, 0, 0, time.UTC)
		if !dash.LastRetrieval.Equal(expectedTime) {
			t.Errorf("Expected LastRetrieval %v, got %v", expectedTime, dash.LastRetrieval)
		}

		// features
		feat := dash.Features
		if feat.Temperature != -2.5 {
			t.Errorf("Expected Temperature=-2.5, got %f", feat.Temperature)
		}
		if feat.Precipitation != 1.2 {
			t.Errorf("Expected Precipitation=1.2, got %f", feat.Precipitation)
		}
		if feat.Capital != "Oslo" {
			t.Errorf("Expected Capital='Oslo', got '%s'", feat.Capital)
		}
		if feat.Coordinates == nil {
			t.Error("Expected non-nil Coordinates")
		} else {
			if feat.Coordinates.Lat != 59.95 {
				t.Errorf("Expected Lat=59.95, got %f", feat.Coordinates.Lat)
			}
			if feat.Coordinates.Lon != 10.75 {
				t.Errorf("Expected Lon=10.75, got %f", feat.Coordinates.Lon)
			}
		}
		if feat.Population != 5370000 {
			t.Errorf("Expected Population=5370000, got %d", feat.Population)
		}
		if feat.Area != 323802 {
			t.Errorf("Expected Area=323802, got %f", feat.Area)
		}
		expectedRates := map[string]float64{"EUR": 0.095, "USD": 0.10}
		if !reflect.DeepEqual(feat.TargetCurrencies, expectedRates) {
			t.Errorf("Mismatch in TargetCurrencies. Expected %v, got %v",
				expectedRates, feat.TargetCurrencies)
		}
	})

	t.Run("JsonRoundTrip", func(t *testing.T) {
		original := Dashboard{
			Country:       "Sweden",
			ISOCode:       "SE",
			LastRetrieval: time.Date(2025, 4, 1, 16, 0, 0, 0, time.UTC),
			Features: DashboardFeatures{
				Temperature:   5.8,
				Precipitation: 0.3,
				Capital:       "Stockholm",
				Coordinates:   &Coordinates{Lat: 59.33, Lon: 18.07},
				Population:    10500000,
				Area:          450295.0,
				TargetCurrencies: map[string]float64{
					"NOK": 0.98, "JPY": 13.3,
				},
			},
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal Dashboard: %v", err)
		}
		var parsed Dashboard
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal Dashboard: %v", err)
		}
		if original.Country != parsed.Country || original.ISOCode != parsed.ISOCode {
			t.Errorf("Mismatch in Country or ISOCode.\nOriginal: %+v\nParsed:   %+v",
				original, parsed)
		}
		if !original.LastRetrieval.Equal(parsed.LastRetrieval) {
			t.Errorf("Mismatch in LastRetrieval times.\nOriginal: %v\nParsed:   %v",
				original.LastRetrieval, parsed.LastRetrieval)
		}
		if !reflect.DeepEqual(original.Features, parsed.Features) {
			t.Errorf("Mismatch in Features.\nOriginal: %+v\nParsed:   %+v",
				original.Features, parsed.Features)
		}
	})

	t.Run("EmptyDashboard", func(t *testing.T) {
		var emptyDash Dashboard
		data, err := json.Marshal(emptyDash)
		if err != nil {
			t.Fatalf("Failed to marshal empty Dashboard: %v", err)
		}
		var parsed Dashboard
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal empty Dashboard: %v", err)
		}
		if !reflect.DeepEqual(emptyDash, parsed) {
			t.Errorf("Empty Dashboard mismatch.\nOriginal: %+v\nParsed:   %+v",
				emptyDash, parsed)
		}
	})
}

func TestCoordinatesStruct(t *testing.T) {
	t.Run("BasicCoordinates", func(t *testing.T) {
		coords := Coordinates{Lat: 34.05, Lon: -118.25}
		if coords.Lat != 34.05 {
			t.Errorf("Expected Lat=34.05, got %f", coords.Lat)
		}
		if coords.Lon != -118.25 {
			t.Errorf("Expected Lon=-118.25, got %f", coords.Lon)
		}
	})

	t.Run("JsonSerialization", func(t *testing.T) {
		coords := Coordinates{Lat: 60.17, Lon: 24.93}
		data, err := json.Marshal(coords)
		if err != nil {
			t.Fatalf("Failed to marshal Coordinates: %v", err)
		}
		var parsed Coordinates
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal Coordinates: %v", err)
		}
		if coords.Lat != parsed.Lat || coords.Lon != parsed.Lon {
			t.Errorf("Coordinates differ.\nOriginal: %+v\nParsed:   %+v",
				coords, parsed)
		}
	})
}
