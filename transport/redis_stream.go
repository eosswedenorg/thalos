
package transport

import (
    "fmt"
    "strings"
    "eosio-ship-trace-reader/redis"
)

type RedisStream struct {
    name string

    // Length of the stream, if items are added when the stream is full, old items will be evicted
    // until the stream's length is equal to this value.
    length int64

    // map of namespaces and their indexes.
    // each namespace is it's own stream.
    indexes map[string]uint32
}

func NewRedisStream(name string, length int64) (RedisStream) {
    return RedisStream{
        name: name,
        length: length,
        indexes: make(map[string]uint32),
    }
}

func (this RedisStream) Send(namespace string, id uint32, message interface{}) error {

    stream := strings.Join([]string{"ship.stream", this.name, namespace}, ".")
    index := this.nextIndex(namespace)

    data := map[string]interface{}{
        "block": id,
        "data": message,
    }

    if err := redis.XAdd(stream, fmt.Sprintf("%d-%d", id, index), this.length, data).Err(); err != nil {
        return fmt.Errorf("Failed to add to redis stream '%s': %s", stream, err)
    }
    return nil
}

func (this RedisStream) Commit() error {

    // reset indexes on flush.
    this.indexes = make(map[string]uint32)
    return nil
}

func (this RedisStream) nextIndex(namespace string) uint32 {
    idx := this.indexes[namespace]
    this.indexes[namespace] = idx + 1
    return idx
}
