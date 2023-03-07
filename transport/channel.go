package transport

import (
	"strings"
)

// Generic interface for all channel types.
type ChannelInterface interface {
	String() string
}

// Standard channel. Just a wrapper around string slice
type Channel []string

func (c *Channel) Append(name ...string) {
	*c = append(*c, name...)
}

func (c Channel) String() string {
	return strings.Join(c, "/")
}

// Check if two channels are equal
func (c Channel) Is(other Channel) bool {
	if len(c) != len(other) {
		return false
	}

	for i, item := range c {
		if item != other[i] {
			return false
		}
	}

	return true
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
		ch.Append("contract", ac.Contract)
	}

	if len(ac.Action) > 0 {
		ch.Append("action", ac.Action)
	}

	return ch.String()
}
