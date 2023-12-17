package cache

import (
	"context"
	"time"
)

type Store interface {
	// Set an item in the store.
	Set(ctx context.Context, key string, value any, TTL time.Duration) error

	// Get an item from the store.
	// returns an error if key is not found or there is other problems.
	Get(ctx context.Context, key string, value any) error

	// Check if a key exist in the store.
	Has(ctx context.Context, key string) bool
}
