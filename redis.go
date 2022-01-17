
package main

import (
    "strings"
    "context"
    "time"
    "github.com/go-redis/redis/v8"
)

var rdb *redis.Client

var redisCtx = context.Background()

var redisPrefix = "ship."

func RedisConnect(addr string, password string, db int) {
    rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
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
