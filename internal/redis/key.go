package redis

import (
	"fmt"
)

// Key consists of a namespace and a channel.
// And is encoded to a string in this format: `<namespace>::<channel>`

type Key struct {
	NS      Namespace
	Channel ChannelInterface
}

func (k Key) String() string {
	return fmt.Sprintf("%s::%s", k.NS, k.Channel)
}
