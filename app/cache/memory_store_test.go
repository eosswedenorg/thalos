package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type memoryTestItem struct {
	String string
	Num    uint32
	Float  float32
}

func TestMemoryStore_Set(t *testing.T) {
	now = func() time.Time { return time.Unix(1581315270, 0) }

	item := memoryTestItem{
		String: "MyString",
		Num:    23,
		Float:  3.14,
	}

	expected := map[string]memoryStoreItem{
		"key1": {
			value:   item,
			expired: now().Add(time.Hour),
		},
	}

	store := NewMemoryStore()
	err := store.Set("key1", item, time.Hour)
	assert.NoError(t, err)

	assert.Equal(t, expected, store.data)
}

func TestMemoryStore_GetMiss(t *testing.T) {
	store := NewMemoryStore()

	var v any
	err := store.Get("Key2", &v)
	assert.Error(t, err)
}

func TestMemoryStore_GetHit(t *testing.T) {
	expected := memoryTestItem{
		String: "MyString",
		Num:    23,
		Float:  3.14,
	}

	store := NewMemoryStore()
	err := store.Set("key1", expected, time.Hour)
	assert.NoError(t, err)

	var actual memoryTestItem
	err = store.Get("key1", &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestMemoryStore_GetNonPointer(t *testing.T) {
	expected := memoryTestItem{
		String: "MyString",
		Num:    23,
		Float:  3.14,
	}

	store := NewMemoryStore()
	err := store.Set("key1", expected, time.Hour)
	assert.NoError(t, err)

	var actual string
	err = store.Get("key1", actual)
	assert.EqualError(t, err, "value must be of pointer type, 'string' passed")
}

func TestMemoryStore_Has(t *testing.T) {
	store := NewMemoryStore()
	err := store.Set("key1", "value", time.Hour)
	assert.NoError(t, err)

	assert.True(t, store.Has("key1"))
	assert.False(t, store.Has("key2"))
}
