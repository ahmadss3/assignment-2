// File: assignment-2/structs/cache_test.go
package structs

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// TestCacheEntryStruct verifies the behavior of the CacheEntry struct
// by assigning sample data, checking fields, and doing JSON round trips.
func TestCacheEntryStruct(t *testing.T) {

	t.Run("FieldAssignments", func(t *testing.T) {
		// Create a CacheEntry with some sample values.
		entry := CacheEntry{
			Key:         "country:NO",
			Data:        []byte(`{"latitude":62.0,"longitude":10.0}`),
			LastFetched: time.Date(2025, 4, 1, 10, 0, 0, 0, time.UTC),
			TTLHours:    24,
		}

		// Check each field for correctness.
		if entry.Key != "country:NO" {
			t.Errorf("Expected Key to be 'country:NO', got '%s'", entry.Key)
		}
		expectedData := []byte(`{"latitude":62.0,"longitude":10.0}`)
		if !reflect.DeepEqual(entry.Data, expectedData) {
			t.Errorf("Data field mismatch.\nExpected: %s\nGot:      %s", expectedData, entry.Data)
		}
		expectedTime := time.Date(2025, 4, 1, 10, 0, 0, 0, time.UTC)
		if !entry.LastFetched.Equal(expectedTime) {
			t.Errorf("Expected LastFetched = %v, got %v", expectedTime, entry.LastFetched)
		}
		if entry.TTLHours != 24 {
			t.Errorf("Expected TTLHours=24, got %d", entry.TTLHours)
		}
	})

	t.Run("ZeroValue", func(t *testing.T) {
		// A zero-value CacheEntry should have empty fields and default values.
		var emptyEntry CacheEntry
		if emptyEntry.Key != "" {
			t.Errorf("Expected Key to be empty, got '%s'", emptyEntry.Key)
		}
		if len(emptyEntry.Data) != 0 {
			t.Errorf("Expected Data to be empty, got something else")
		}
		if !emptyEntry.LastFetched.IsZero() {
			t.Errorf("Expected LastFetched to be zero time, got %v", emptyEntry.LastFetched)
		}
		if emptyEntry.TTLHours != 0 {
			t.Errorf("Expected TTLHours=0, got %d", emptyEntry.TTLHours)
		}
	})

	t.Run("JsonRoundTrip", func(t *testing.T) {
		original := CacheEntry{
			Key:         "meteo:59.33,18.07",
			Data:        []byte(`{"hourly":[5.5, 6.2, 7.0],"timestamp":"2025-04-01T09:00:00Z"}`),
			LastFetched: time.Date(2025, 4, 1, 9, 0, 0, 0, time.UTC),
			TTLHours:    1,
		}

		// Marshal to JSON.
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal CacheEntry: %v", err)
		}

		// Unmarshal into another instance.
		var parsed CacheEntry
		if err := json.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal CacheEntry: %v", err)
		}

		// Check equality.
		if original.Key != parsed.Key {
			t.Errorf("Key mismatch after JSON round trip.\nOriginal: %s\nParsed:   %s", original.Key, parsed.Key)
		}
		if !reflect.DeepEqual(original.Data, parsed.Data) {
			t.Errorf("Data mismatch after JSON round trip.\nOriginal: %s\nParsed:   %s", original.Data, parsed.Data)
		}
		if !original.LastFetched.Equal(parsed.LastFetched) {
			t.Errorf("LastFetched mismatch.\nOriginal: %v\nParsed:   %v",
				original.LastFetched, parsed.LastFetched)
		}
		if original.TTLHours != parsed.TTLHours {
			t.Errorf("TTLHours mismatch.\nOriginal: %d\nParsed:   %d",
				original.TTLHours, parsed.TTLHours)
		}
	})
}
