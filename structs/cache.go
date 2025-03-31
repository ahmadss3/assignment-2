// File: assignment-2/structs/cache.go
package structs

import "time"

// CacheEntry stores a piece of data with a key, the data itself in a byte slice,
// the time at which it was last fetched, and a TTL (in hours) indicating how long
// the data should remain valid.
type CacheEntry struct {
	Key         string    `firestore:"key"`         // Key is the unique identifier for this cache record.
	Data        []byte    `firestore:"data"`        // Data is a byte array that can hold any serialized payload.
	LastFetched time.Time `firestore:"lastFetched"` // LastFetched indicates the point in time at which this data was obtained from an external source.
	TTLHours    int       `firestore:"ttlHours"`    // TTLHours specifies how many hours this data should be considered valid before it is purged or re-fetched.

}
