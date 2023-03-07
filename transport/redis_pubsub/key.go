package redis_pubsub

import (
	"fmt"

	"eosio-ship-trace-reader/transport"
)

// Key consists of a namespace and a channel.
// And is encoded to a string in this format: `<namespace>::<channel>`

type Key struct {
	NS      Namespace
	Channel transport.Channel
}

func (k Key) String() string {
	return fmt.Sprintf("%s::%s", k.NS, k.Channel)
}
