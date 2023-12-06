package cache

import "time"

type Store interface {
	// Set an item in the store.
	Set(key string, value any, TTL time.Duration) error

	// Get an item from the store.
	// returns an error if key is not found or there is other problems.
	Get(key string, value any) error

	// Check if a key exist in the store.
	Has(key string) bool
}
