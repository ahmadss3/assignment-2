// File: assignment-2/firebase/cache_firebase.go
package firebase

import (
	"context"
	"fmt"
	"time"

	"assignment-2/constants"
	"assignment-2/structs"
)

// GetCacheEntry fetches a cache document by 'key' from the Firestore "cache" collection.
// If the document does not exist, an error is returned. If data cannot be parsed, an error is returned.
func GetCacheEntry(ctx context.Context, key string) (*structs.CacheEntry, error) {
	if err := ensureClient(); err != nil {
		return nil, err
	}

	docRef := FirestoreClient.Collection(constants.CACHE_COLLECTION).Doc(key)
	snap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache doc for key=%s: %v", key, err)
	}
	if !snap.Exists() {
		return nil, fmt.Errorf("cache doc not found for key=%s", key)
	}

	var ce structs.CacheEntry
	if err := snap.DataTo(&ce); err != nil {
		return nil, fmt.Errorf("failed to parse cache doc: %v", err)
	}
	return &ce, nil
}

// SaveCacheEntry creates or updates a cache document in the Firestore "cache" collection.
// Requires 'key' plus 'data', 'lastFetched', and 'ttlHours' fields in the CacheEntry struct.
func SaveCacheEntry(ctx context.Context, entry structs.CacheEntry) error {
	if err := ensureClient(); err != nil {
		return err
	}

	docRef := FirestoreClient.Collection(constants.CACHE_COLLECTION).Doc(entry.Key)
	_, err := docRef.Set(ctx, map[string]interface{}{
		"key":         entry.Key,
		"data":        entry.Data,
		"lastFetched": entry.LastFetched,
		"ttlHours":    entry.TTLHours,
	})
	if err != nil {
		return fmt.Errorf("failed to save cache doc (key=%s): %v", entry.Key, err)
	}
	return nil
}

// PurgeOldCache deletes cache docs older than the given duration from the Firestore "cache" collection.
// It compares the "lastFetched" timestamp to (time.Now() - olderThan) and removes stale entries.
func PurgeOldCache(ctx context.Context, olderThan time.Duration) error {
	if err := ensureClient(); err != nil {
		return err
	}

	cutoff := time.Now().Add(-olderThan)
	colRef := FirestoreClient.Collection(constants.CACHE_COLLECTION)
	q := colRef.Where("lastFetched", "<", cutoff)
	snaps, err := q.Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("failed to query old cache docs: %v", err)
	}

	for _, s := range snaps {
		_, delErr := s.Ref.Delete(ctx)
		if delErr != nil {
			fmt.Printf("Warning: failed to delete old cache doc %s: %v\n", s.Ref.ID, delErr)
		}
	}
	return nil
}
