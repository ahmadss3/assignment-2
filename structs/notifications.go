// File: assignment-2/structs/notifications.go
package structs

import "time"

// Notification represents a single webhook registration that is triggered when certain events occur.
type Notification struct {
	ID      string    `json:"id,omitempty"`      // ID is the unique identifier for this webhook registration.
	URL     string    `json:"url"`               // Is the destination for the POST request when an event matching this webhook is triggered.
	Country string    `json:"country,omitempty"` // Country indicates the country filter. If empty, the webhook applies to all countries.
	Event   string    `json:"event"`             // Event is the type of event on which the webhook triggers ("REGISTER", "CHANGE", "DELETE", "INVOKE").
	Created time.Time `json:"created,omitempty"` // Created is the time at which this webhook registration was initially created.
}
