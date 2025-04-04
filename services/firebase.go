// File: assignment-2/services/firebase.go
package services

import (
	"context"
	"errors"

	"assignment-2/structs"
)

// Define variables (function pointers) to allow test patching.
var (
	GetCacheEntryVar  = getCacheEntry
	SaveCacheEntryVar = saveCacheEntry
)

func getCacheEntry(ctx context.Context, key string) (*structs.CacheEntryEx, error) {
	// Return nil "not found"
	return nil, errors.New("stub: no cache entry found")
}

func saveCacheEntry(ctx context.Context, entry structs.CacheEntryEx) error {
	// Do nothing
	return nil
}
