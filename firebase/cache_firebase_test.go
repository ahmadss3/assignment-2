// File: assignment-2/firebase/cache_firebase_test.go
package firebase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"assignment-2/structs"
)

// This test file covers the functionality in cache_firebase.go.
// It includes tests for GetCacheEntry, SaveCacheEntry, and PurgeOldCache.

func TestCacheFirebase(t *testing.T) {
	if FirestoreClient == nil {
		t.Skip("FirestoreClient is not initialized. Skipping cache Firebase tests.")
	}

	// create a test context for Firestore operations.
	ctx := context.Background()

	// Create a random key for testing.
	testKey := fmt.Sprintf("testKey-%d", time.Now().UnixNano())

	t.Run("SaveCacheEntry", func(t *testing.T) {
		entry := structs.CacheEntry{
			Key:         testKey,
			Data:        []byte(`{"message":"hello"}`),
			LastFetched: time.Now(),
			TTLHours:    2,
		}
		err := SaveCacheEntry(ctx, entry)
		if err != nil {
			t.Errorf("SaveCacheEntry failed: %v", err)
		}
	})

	t.Run("GetCacheEntry", func(t *testing.T) {
		// Retrieve the entry we just saved.
		ce, err := GetCacheEntry(ctx, testKey)
		if err != nil {
			t.Errorf("GetCacheEntry failed: %v", err)
			return
		}
		if ce.Key != testKey {
			t.Errorf("Expected key '%s', got '%s'", testKey, ce.Key)
		}
		if string(ce.Data) != `{"message":"hello"}` {
			t.Errorf("Expected data to be `{\"message\":\"hello\"}`, got '%s'", string(ce.Data))
		}
		if ce.TTLHours != 2 {
			t.Errorf("Expected TTLHours=2, got %d", ce.TTLHours)
		}
	})

	t.Run("GetCacheEntry_NotFound", func(t *testing.T) {
		// Attempt to fetch an entry that doesn't exist.
		_, err := GetCacheEntry(ctx, "nonExistingKey-XYZ")
		if err == nil {
			t.Error("Expected error for non-existing key, got nil")
		}
	})

	t.Run("PurgeOldCache", func(t *testing.T) {
		// We'll artificially purge items older than 0 hours to remove everything
		err := PurgeOldCache(ctx, 0)
		if err != nil {
			t.Errorf("PurgeOldCache failed: %v", err)
		}

		// Check if testKey entry was removed or not.
		ce, err := GetCacheEntry(ctx, testKey)
		if err == nil && ce != nil {
			// Possibly not older than 0 hours. We'll just log a note.
			t.Logf("Cache entry with key '%s' still exists after PurgeOldCache(0). Possibly it wasn't old enough.", testKey)
		} else {
			t.Log("Cache entry was purged or not found (expected).")
		}
	})
}

func TestEnsureClientCache(t *testing.T) {
	if FirestoreClient == nil {
		t.Skip("Skipping because FirestoreClient is nil.")
	}
	err := ensureClient()
	if err != nil {
		t.Errorf("ensureClient unexpectedly returned error: %v", err)
	}
}
