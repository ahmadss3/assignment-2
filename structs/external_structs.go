package structs

import "time"

// CacheEntryEx holds cached data, typically stored in Firestore for external APIs.
type CacheEntryEx struct {
	Key         string    `firestore:"key"`
	Data        []byte    `firestore:"data"`
	LastFetched time.Time `firestore:"lastFetched"`
	TTLHours    int       `firestore:"ttlHours"`
}

// CountryInfo holds data about a country for external usage.
type CountryInfo struct {
	Name         string
	Capital      string
	Population   int64
	Area         float64
	BaseCurrency string
	Coordinates  Coordinates
}

// Coordinates represents latitude/longitude
type Coordinates struct {
	Lat float64
	Lon float64
}

// MeteoData holds average temperature and precipitation.
type MeteoData struct {
	AverageTemp          float64
	AveragePrecipitation float64
}

// CurrencyRates is a map from currency code to exchange rate.
type CurrencyRates map[string]float64

// Dashboard represents the data returned by GET /dashboard/v1/dashboards/{id}.
type Dashboard struct {
	Country       string            `json:"country,omitempty"`
	ISOCode       string            `json:"isoCode,omitempty"`
	Features      DashboardFeatures `json:"features,omitempty"`
	LastRetrieval time.Time         `json:"lastRetrieval,omitempty"`
}

// DashboardFeatures contains the data that the user requested for inclusion in the dashboard.
// It merges the old "DashboardFeatures" definitions with optional fields.
type DashboardFeatures struct {
	Temperature      float64            `json:"temperature,omitempty"`
	Precipitation    float64            `json:"precipitation,omitempty"`
	Capital          string             `json:"capital,omitempty"`
	Coordinates      *Coordinates       `json:"coordinates,omitempty"`
	Population       int64              `json:"population,omitempty"`
	Area             float64            `json:"area,omitempty"`
	TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"`
}
