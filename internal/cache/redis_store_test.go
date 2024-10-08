package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/karlseguin/typed"

	redis_cache "github.com/go-redis/cache/v9"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	Num  uint32
	Name string
}

func TestRedisStore_getOptionsDefaults(t *testing.T) {
	opts := typed.Typed{}

	expected := options{
		Stats: false,
		Size:  1000,
		TTL:   10 * time.Minute,
	}

	actual := getOptions(opts)

	assert.Equal(t, expected, actual)
}

func TestRedisStore_getOptions(t *testing.T) {
	opts := typed.Typed{
		"stats": true,
		"size":  123,
		"ttl":   60,
	}

	expected := options{
		Stats: true,
		Size:  123,
		TTL:   60 * time.Minute,
	}

	actual := getOptions(opts)

	assert.Equal(t, expected, actual)
}

func TestRedisStore_Set(t *testing.T) {
	client, mock := redismock.NewClientMock()

	store := NewRedisStore(&redis_cache.Options{
		Redis: client,
	})

	expected := testItem{
		Num:  24,
		Name: "Some Name",
	}

	bytes, err := store.c.Marshal(expected)
	assert.NoError(t, err)

	mock.ExpectSet("mykey", bytes, time.Minute).SetVal("OK")

	err = store.Set(context.Background(), "mykey", expected, time.Minute)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_GetMiss(t *testing.T) {
	client, mock := redismock.NewClientMock()

	store := NewRedisStore(&redis_cache.Options{
		Redis: client,
	})

	mock.ExpectGet("mykey").SetErr(redis_cache.ErrCacheMiss)

	expected := testItem{}
	err := store.Get(context.Background(), "mykey", &expected)
	assert.ErrorIs(t, err, redis_cache.ErrCacheMiss)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_GetHit(t *testing.T) {
	client, mock := redismock.NewClientMock()

	store := NewRedisStore(&redis_cache.Options{
		Redis: client,
	})

	expected := testItem{
		Num:  42,
		Name: "MyName",
	}

	bytes, err := store.c.Marshal(expected)
	assert.NoError(t, err)

	mock.ExpectSet("mykey2", bytes, time.Second*20).SetVal("OK")
	mock.ExpectGet("mykey2").SetVal(string(bytes))

	err = store.Set(context.Background(), "mykey2", expected, time.Second*20)
	assert.NoError(t, err)

	actual := testItem{}
	err = store.Get(context.Background(), "mykey2", &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStore_Has(t *testing.T) {
	client, mock := redismock.NewClientMock()

	store := NewRedisStore(&redis_cache.Options{
		Redis: client,
	})

	bytes, err := store.c.Marshal("value")
	assert.NoError(t, err)

	mock.ExpectSet("key1", bytes, time.Minute*15).SetVal("OK")
	mock.ExpectGet("key1").SetVal(string(bytes))
	mock.ExpectGet("key2").RedisNil()

	err = store.Set(context.Background(), "key1", "value", time.Minute*15)
	assert.NoError(t, err)
	assert.True(t, store.Has(context.Background(), "key1"))
	assert.False(t, store.Has(context.Background(), "key2"))
	assert.NoError(t, mock.ExpectationsWereMet())
}
