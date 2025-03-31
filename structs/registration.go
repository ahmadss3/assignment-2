// File: assignment-2/structs/registration.go
package structs

import "time"

// Registration describes a configuration used for building a dynamic dashboard.
// This struct stores essential information about the country, its ISO code,
// feature flags for which data to retrieve, and a timestamp for the last change.
type Registration struct {
	// ID is the unique identifier, typically assigned by the database layer.
	ID string `json:"id,omitempty"`

	// Country is the name of the country for which the dashboard is built.
	// Either Country or ISOCode can be used to identify the location.
	Country string `json:"country,omitempty"`

	// ISOCode is the ISO country code (e.g. "NO", "SE", "DE").
	// It can be used in place of the Country field.
	ISOCode string `json:"isoCode,omitempty"`

	// Features holds the boolean flags and target currencies that specify
	// which data to include in the final dashboard (temperature, precipitation, etc.).
	Features Features `json:"features"`

	// LastChange indicates when this registration was last updated on the server side.
	LastChange time.Time `json:"lastChange"`
}
