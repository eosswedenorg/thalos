
package transport

import "fmt"

func Make(driver string, name string) (Driver, error) {

    switch driver {
    case "redis-pubsub":
        return NewRedisPubSub(name), nil
    case "redis-stream":
        return NewRedisStream(name, 1000), nil
    default:
        return nil, fmt.Errorf("Invalid type: %s", driver)
    }
}
