// File: assignment-2/structs/registration_test.go
package structs

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// TestRegistrationStruct verifies basic usage of the Registration struct.
// It checks creation, field assignments, and JSON (un)marshalling.
func TestRegistrationStruct(t *testing.T) {
	t.Run("CreateAndVerifyFields", func(t *testing.T) {
		// Creating an instance of Registration with some test data.
		reg := Registration{
			ID:         "test123",
			Country:    "Sweden",
			ISOCode:    "SE",
			Features:   Features{Temperature: true, Capital: true},
			LastChange: time.Date(2025, 3, 30, 14, 0, 0, 0, time.UTC),
		}

		// Verify field values directly.
		if reg.ID != "test123" {
			t.Errorf("Expected ID to be 'test123', got '%s'", reg.ID)
		}
		if reg.Country != "Sweden" {
			t.Errorf("Expected Country to be 'Sweden', got '%s'", reg.Country)
		}
		if reg.ISOCode != "SE" {
			t.Errorf("Expected ISOCode to be 'SE', got '%s'", reg.ISOCode)
		}
		if !reg.Features.Temperature {
			t.Error("Expected Temperature feature to be true")
		}
		if !reg.Features.Capital {
			t.Error("Expected Capital feature to be true")
		}
		if reg.Features.Precipitation {
			t.Error("Did not expect Precipitation to be set to true")
		}

		// Checking the LastChange timestamp.
		expectedTime := time.Date(2025, 3, 30, 14, 0, 0, 0, time.UTC)
		if !reg.LastChange.Equal(expectedTime) {
			t.Errorf("Expected LastChange to be %v, got %v",
				expectedTime, reg.LastChange)
		}
	})

	t.Run("JsonMarshalUnmarshal", func(t *testing.T) {
		// Construct a sample Registration.
		original := Registration{
			ID:         "abc123",
			Country:    "Norway",
			ISOCode:    "NO",
			Features:   Features{Temperature: true, Population: true},
			LastChange: time.Date(2025, 3, 30, 10, 30, 0, 0, time.UTC),
		}

		// Marshal to JSON.
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal Registration: %v", err)
		}

		// Unmarshal back into a new instance.
		var parsed Registration
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal Registration: %v", err)
		}

		// Compare original and parsed. For slices or nested structs,
		// a deep equality check is useful. Here, we can rely on reflect.DeepEqual.
		if !reflect.DeepEqual(original.ID, parsed.ID) ||
			!reflect.DeepEqual(original.Country, parsed.Country) ||
			!reflect.DeepEqual(original.ISOCode, parsed.ISOCode) {
			t.Errorf("Basic string fields differ after JSON round trip.\nOriginal: %+v\nParsed: %+v", original, parsed)
		}

		// Check the Features struct.
		if !reflect.DeepEqual(original.Features, parsed.Features) {
			t.Errorf("Features differ after JSON round trip.\nOriginal: %+v\nParsed: %+v", original.Features, parsed.Features)
		}

		// Check the LastChange field. We expect the timestamp to match exactly.
		if !original.LastChange.Equal(parsed.LastChange) {
			t.Errorf("LastChange timestamp differs.\nOriginal: %v\nParsed: %v",
				original.LastChange, parsed.LastChange)
		}
	})
}
