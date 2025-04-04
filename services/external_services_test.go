// file assignment-2/services/external_services_test.go
package services

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"assignment-2/structs"
)

// Overridable references
var originalFetchCountryInfo = FetchCountryInfo

func TestFetchCountryInfoStub(t *testing.T) {
	// Demonstraton of stubbing
	FetchCountryInfo = func(countryOrISO string) (*structs.CountryInfo, error) {
		if countryOrISO == "NO" {
			return &structs.CountryInfo{
				Name:         "Norway",
				Capital:      "Oslo",
				Population:   5372000,
				Area:         385207.0,
				BaseCurrency: "NOK",
				Coordinates:  structs.Coordinates{Lat: 60.0, Lon: 10.0},
			}, nil
		}
		return nil, errors.New("stub: only 'NO' supported")
	}
	defer func() { FetchCountryInfo = originalFetchCountryInfo }()

	info, err := FetchCountryInfo("NO")
	if err != nil {
		t.Fatalf("FetchCountryInfo stub error: %v", err)
	}
	if info.Capital != "Oslo" {
		t.Errorf("Expected capital=Oslo, got %s", info.Capital)
	}
}

func TestFetchMeteoDataStub(t *testing.T) {
	orig := FetchMeteoData
	FetchMeteoData = func(lat, lon float64) (*structs.MeteoData, error) {
		return &structs.MeteoData{AverageTemp: 5.5, AveragePrecipitation: 1.2}, nil
	}
	defer func() { FetchMeteoData = orig }()

	data, err := FetchMeteoData(59.0, 10.0)
	if err != nil {
		t.Fatalf("FetchMeteoData stub error: %v", err)
	}
	if data.AverageTemp != 5.5 {
		t.Errorf("Expected 5.5, got %f", data.AverageTemp)
	}
}

func TestFetchCurrencyRatesStub(t *testing.T) {
	orig := FetchCurrencyRates
	FetchCurrencyRates = func(base string) (structs.CurrencyRates, error) {
		if base == "NOK" {
			return structs.CurrencyRates{"EUR": 0.09, "USD": 0.1}, nil
		}
		return nil, errors.New("stub: only 'NOK' supported")
	}
	defer func() { FetchCurrencyRates = orig }()

	rates, err := FetchCurrencyRates("NOK")
	if err != nil {
		t.Fatalf("FetchCurrencyRates stub error: %v", err)
	}
	if val, ok := rates["EUR"]; !ok || val != 0.09 {
		t.Errorf("Expected EUR=0.09, got %f or missing", val)
	}
}

func TestMockDataIntegration(t *testing.T) {
	mockFile := "mock_files/restcountries_norway.json"
	if _, err := os.Stat(mockFile); os.IsNotExist(err) {
		t.Skipf("Skipping because mock file not found: %s", mockFile)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile(mockFile)
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("mock restcountries GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("mock restcountries returned %d => %s", resp.StatusCode, body)
	}
	t.Log("Successfully read from local mock file in test.")
}
