// File: assignment-2/structs/dashboard.go
package structs

import "time"

// Dashboard represents the data returned by GET /dashboard/v1/dashboards/{id}.
// It aggregates various details such as Country, ISOCode, a set of requested Features,
// and the time at which the dashboard data was last retrieved.

type Dashboard struct {
	Country       string            `json:"country,omitempty"`       // The name of the country.
	ISOCode       string            `json:"isoCode,omitempty"`       // ISOCode is the two-letter country code.
	Features      DashboardFeatures `json:"features,omitempty"`      // Features holds the details about temperature, precipitation, capital, coordinates, etc.// The presence of each field depends on user-selected options
	LastRetrieval time.Time         `json:"lastRetrieval,omitempty"` // The timestamp at which the dashboard data was last fetched.
}

// DashboardFeatures contains the data that the user requested for inclusion in the dashboard.

type DashboardFeatures struct {
	Temperature      float64            `json:"temperature,omitempty"`      // Temperature is the value.
	Precipitation    float64            `json:"precipitation,omitempty"`    // Precipitation is the value.
	Capital          string             `json:"capital,omitempty"`          // The country's capital city.
	Coordinates      *Coordinates       `json:"coordinates,omitempty"`      // Coordinates is an optional pointer to geographical coordinates (latitude, longitude).If the user has not requested coordinates, it can be nil.
	Population       int64              `json:"population,omitempty"`       // Population is the total population of the country.
	Area             float64            `json:"area,omitempty"`             // Area is the total land area for the country.
	TargetCurrencies map[string]float64 `json:"targetCurrencies,omitempty"` // TargetCurrencies holds a map of currency codes to their exchange rates, relative to the country's base currency.
}

// Coordinates stores latitude and longitude information for the relevant country or region.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`  // Is the geographic latitude.
	Longitude float64 `json:"longitude"` // Is the geographic longitude.
}
