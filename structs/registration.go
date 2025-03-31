// File: assignment-2/structs/registration.go
package structs

import "time"

// Registration describes a configuration used for building a dashboard.

type Registration struct {
	ID         string    `json:"id,omitempty"`      // ID is the unique identifier
	Country    string    `json:"country,omitempty"` // Country is the name.
	ISOCode    string    `json:"isoCode,omitempty"` // ISOCode is the ISO country code.
	Features   Features  `json:"features"`          // Features holds the boolean flags and target currencies that specify.
	LastChange time.Time `json:"lastChange"`        // LastChange indicates when this registration was last updated.
}
