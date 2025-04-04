// File: assignment-2/services/firebase_test.go
package services

import (
	"context"
	"errors"
	"testing"

	"assignment-2/structs"
)

// TestFirebaseServiceStubs demonstrate use og GetCacheEntryVar og SaveCacheEntryVar
// witch they er stubbed in firebase.go.
func TestFirebaseServiceStubs(t *testing.T) {
	t.Run("TestGetCacheEntryStub", func(t *testing.T) {
		// call GetCacheEntryVar witch points to getCacheEntry
		entry, err := GetCacheEntryVar(context.Background(), "test-key")
		if entry != nil {
			t.Error("Expected entry to be nil from the stub, but got a non-nil value")
		}
		if err == nil {
			t.Error("Expected an error from the stub, but got nil")
		} else {
			// Check that the error matches the stub return value
			expectedErr := "stub: no cache entry found"
			if err.Error() != expectedErr {
				t.Errorf("Expected error=%q, got %q", expectedErr, err.Error())
			}
		}
	})

	t.Run("TestSaveCacheEntryStub", func(t *testing.T) {
		// Calls SaveCacheEntryVar which currently points to saveCacheEntry
		testEntry := structs.CacheEntryEx{
			Key: "test-key",
		}
		err := SaveCacheEntryVar(context.Background(), testEntry)
		if err != nil {
			t.Errorf("Expected no error from the stub, got %v", err)
		}
	})

	t.Run("TestOverrideStubs", func(t *testing.T) {
		originalGet := GetCacheEntryVar
		originalSave := SaveCacheEntryVar
		defer func() {
			// Revert patching after this test
			GetCacheEntryVar = originalGet
			SaveCacheEntryVar = originalSave
		}()

		// Override GetCacheEntryVar
		GetCacheEntryVar = func(ctx context.Context, key string) (*structs.CacheEntryEx, error) {
			return &structs.CacheEntryEx{Key: key}, nil
		}
		// Override SaveCacheEntryVar
		SaveCacheEntryVar = func(ctx context.Context, entry structs.CacheEntryEx) error {
			if entry.Key == "" {
				return errors.New("test: empty key is invalid")
			}
			return nil
		}

		e, err := GetCacheEntryVar(context.Background(), "patched-key")
		if err != nil {
			t.Errorf("Expected no error from the patched function, got %v", err)
		}
		if e == nil || e.Key != "patched-key" {
			t.Errorf("Expected e.Key='patched-key', got %v", e)
		}

		// Test of the patched SaveCacheEntryVar
		err = SaveCacheEntryVar(context.Background(), structs.CacheEntryEx{Key: ""})
		if err == nil {
			t.Error("Expected error for empty key, got nil")
		}
	})
}
