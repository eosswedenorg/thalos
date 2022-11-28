
package transport

import (
    "fmt"
    "strings"
    "eosio-ship-trace-reader/redis"
)

type RedisPubSub struct {
    name string
}

func NewRedisPubSub(name string) (RedisPubSub) {
    return RedisPubSub{
        name: name,
    }
}

func (this RedisPubSub) Send(namespace string, id uint32, message interface{}) error {

    channel := strings.Join([]string{"ship.channel", this.name, namespace}, ".")
    if err := redis.RegisterPublish(channel, message).Err(); err != nil {
        return fmt.Errorf("Failed to post to channel '%s': %s", channel, err)
    }
    return nil
}

func (this RedisPubSub) Commit() error {
    _, err := redis.Send()
    if err != nil {
        return fmt.Errorf("Failed to send redis. command: %s", err)
    }
    return nil
}
