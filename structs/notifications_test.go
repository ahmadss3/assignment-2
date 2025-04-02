// File: assignment-2/structs/notifications_test.go
package structs

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// TestNotificationStruct covers basic usage of the Notification struct,
// including field assignments and JSON marshalling/unmarshalling.
func TestNotificationStruct(t *testing.T) {

	t.Run("FieldAssignments", func(t *testing.T) {
		// Create a Notification instance with sample data
		notif := Notification{
			ID:      "notif123",
			URL:     "https://webhook.example.com/endpoint",
			Country: "NO",
			Event:   "REGISTER",
			Created: time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC),
		}

		// Verify fields individually
		if notif.ID != "notif123" {
			t.Errorf("Expected ID to be 'notif123', got '%s'", notif.ID)
		}
		if notif.URL != "https://webhook.example.com/endpoint" {
			t.Errorf("Expected URL to be 'https://webhook.example.com/endpoint', got '%s'", notif.URL)
		}
		if notif.Country != "NO" {
			t.Errorf("Expected Country to be 'NO', got '%s'", notif.Country)
		}
		if notif.Event != "REGISTER" {
			t.Errorf("Expected Event to be 'REGISTER', got '%s'", notif.Event)
		}

		expectedCreated := time.Date(2025, 4, 1, 12, 0, 0, 0, time.UTC)
		if !notif.Created.Equal(expectedCreated) {
			t.Errorf("Expected Created to be %v, got %v", expectedCreated, notif.Created)
		}
	})

	t.Run("JsonRoundTrip", func(t *testing.T) {
		// Construct an instance with some data
		original := Notification{
			ID:      "abcd-789",
			URL:     "https://mywebhook.com/listener",
			Country: "",
			Event:   "CHANGE",
			Created: time.Date(2025, 3, 30, 18, 30, 0, 0, time.UTC),
		}

		// Marshal to JSON
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal Notification: %v", err)
		}

		// Unmarshal into a new instance
		var parsed Notification
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal Notification: %v", err)
		}

		// Compare original and parsed.
		// For simpler fields, direct equality is enough. For nested or complex fields, use reflect.DeepEqual.
		if original.ID != parsed.ID ||
			original.URL != parsed.URL ||
			original.Country != parsed.Country ||
			original.Event != parsed.Event {
			t.Errorf("One or more basic fields differ after JSON round trip.\nOriginal: %+v\nParsed:   %+v",
				original, parsed)
		}

		if !original.Created.Equal(parsed.Created) {
			t.Errorf("Created time differs.\nOriginal: %v\nParsed:   %v",
				original.Created, parsed.Created)
		}
	})

	t.Run("EmptyFields", func(t *testing.T) {
		// A scenario where certain fields are left empty
		emptyNotif := Notification{
			URL:   "https://webhook.empty-country.com",
			Event: "DELETE",
		}

		// Marshal
		data, err := json.Marshal(emptyNotif)
		if err != nil {
			t.Fatalf("Failed to marshal Notification with empty fields: %v", err)
		}

		// Unmarshal
		var parsed Notification
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal Notification with empty fields: %v", err)
		}

		// Use reflect.DeepEqual for the entire struct if desired
		if !reflect.DeepEqual(emptyNotif, parsed) {
			t.Errorf("Notification struct differs after JSON round trip.\nOriginal: %+v\nParsed:   %+v",
				emptyNotif, parsed)
		}
	})
}
