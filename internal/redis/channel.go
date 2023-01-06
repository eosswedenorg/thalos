package redis

import (
	"fmt"
	"strings"
)

// Generic interface for all channel types.
type ChannelInterface interface {
	String() string
}

// Standard channel. Just a wrapper around string slice
type Channel []string

func (c *Channel) Append(name string) {
	*c = append(*c, name)
}

func (c Channel) String() string {
	return strings.Join(c, ".")
}

// Define channels without any variables.
var (
	TransactionChannel = Channel{"transaction"}
	HeartbeatChannel   = Channel{"heartbeat"}
)

// Action channel.

type ActionChannel struct {
	Contract string
	Action   string
}

func (ac ActionChannel) String() string {
	ch := Channel{"actions"}

	if len(ac.Contract) > 0 {
		ch.Append(fmt.Sprintf("contract:%s", ac.Contract))
	}

	if len(ac.Action) > 0 {
		ch.Append(fmt.Sprintf("action:%s", ac.Action))
	}

	return ch.String()
}
