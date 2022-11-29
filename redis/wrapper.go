package redis

import (
	"context"
	"strings"
	"time"

	_redis "github.com/go-redis/redis/v8"
)

var rdb *_redis.Client

var redis_pipe _redis.Pipeliner

var redisCtx = context.Background()

var prefix []string

func Connect(addr string, password string, db int) error {
	rdb = _redis.NewClient(&_redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	redis_pipe = rdb.Pipeline()

	return rdb.Ping(redisCtx).Err()
}

func Client() *_redis.Client {
	return rdb
}

func SetPrefix(components ...string) {
	prefix = components
}

func Key(components ...string) string {
	components = append(prefix, components...)
	return strings.Join(components, ".")
}

func Get(key string) *_redis.StringCmd {
	return rdb.Get(redisCtx, key)
}

func Set(key string, value interface{}, expiration time.Duration) *_redis.StatusCmd {
	return rdb.Set(redisCtx, key, value, expiration)
}

func Publish(channel string, message interface{}) *_redis.IntCmd {
	return rdb.Publish(redisCtx, channel, message)
}

func RegisterPublish(channel string, message interface{}) *_redis.IntCmd {
	return redis_pipe.Publish(redisCtx, channel, message)
}

func Send() ([]_redis.Cmder, error) {
	return redis_pipe.Exec(redisCtx)
}
