package cache

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

// Store time function in a variable.
// Makes it easy to travel in time when testing.
var now = time.Now

type memoryStoreItem struct {
	// Actual value stored.
	value any

	// Cache expiration time.
	expired time.Time
}

type MemoryStore struct {
	data map[string]memoryStoreItem
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{make(map[string]memoryStoreItem)}
}

func (s *MemoryStore) Get(ctx context.Context, key string, value any) error {
	if item, ok := s.data[key]; ok {

		if item.expired.Before(now()) {
			delete(s.data, key)
			return fmt.Errorf("key: %s does not exist", key)
		}

		v := reflect.ValueOf(value)
		if v.Kind() != reflect.Pointer {
			return fmt.Errorf("value must be of pointer type, '%s' passed", v.Kind().String())
		}

		v.Elem().Set(reflect.ValueOf(item.value))

		return nil
	}
	return fmt.Errorf("key: %s does not exist", key)
}

func (s *MemoryStore) Has(ctx context.Context, key string) bool {
	_, hit := s.data[key]
	return hit
}

func (s *MemoryStore) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	s.data[key] = memoryStoreItem{
		value:   value,
		expired: now().Add(ttl),
	}
	return nil
}
