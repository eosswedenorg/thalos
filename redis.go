
package main

import (
    "strings"
    "context"
    "time"
    "github.com/go-redis/redis/v8"
)

var rdb *redis.Client

var redis_pipe redis.Pipeliner

var redisCtx = context.Background()

var redisPrefix = "ship."

func RedisConnect(addr string, password string, db int) error {
    rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    redis_pipe = rdb.Pipeline()

    return rdb.Ping(redisCtx).Err()
}

func RedisKey(components ...string) (string) {
    return redisPrefix + strings.Join(components, ".")
}

func RedisGet(key string) (*redis.StringCmd) {
    return rdb.Get(redisCtx, key)
}

func RedisSet(key string, value interface{}, expiration time.Duration) (*redis.StatusCmd) {
    return rdb.Set(redisCtx, key, value, expiration)
}

func RedisPublish(channel string, message interface{}) (*redis.IntCmd) {
    return rdb.Publish(redisCtx, channel, message)
}

func RedisRegisterPublish(channel string, message interface{}) (*redis.IntCmd) {
    return redis_pipe.Publish(redisCtx, channel, message)
}

func RedisSend() ([]redis.Cmder, error) {
    return redis_pipe.Exec(redisCtx)
}
