
package redis

import (
    "strings"
    "context"
    "time"
    _redis "github.com/go-redis/redis/v8"
)

var rdb *_redis.Client

var redis_pipe _redis.Pipeliner

var redisCtx = context.Background()

var Prefix = "ship."

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

func Key(components ...string) (string) {
    return Prefix + strings.Join(components, ".")
}

func Get(key string) (*_redis.StringCmd) {
    return rdb.Get(redisCtx, key)
}

func Set(key string, value interface{}, expiration time.Duration) (*_redis.StatusCmd) {
    return rdb.Set(redisCtx, key, value, expiration)
}

func Publish(channel string, message interface{}) (*_redis.IntCmd) {
    return rdb.Publish(redisCtx, channel, message)
}

func XAdd(stream string, id string, maxlen int64, message map[string]interface{}) (*_redis.StringCmd) {

    args := &_redis.XAddArgs{
    	Stream: stream,
    	ID: id,
        MaxLenApprox: maxlen,
    	Values: message,
    }

    return rdb.XAdd(redisCtx, args)
}

func RegisterPublish(channel string, message interface{}) (*_redis.IntCmd) {
    return redis_pipe.Publish(redisCtx, channel, message)
}

func Send() ([]_redis.Cmder, error) {
    return redis_pipe.Exec(redisCtx)
}
