package transport

import (
	"strings"
)

// Channel is just a wrapper around string slice
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

// Action Channel
type Action struct {
	Contract string
	Action   string
}

func (a Action) Channel() Channel {
	ch := Channel{"actions"}

	if len(a.Contract) > 0 {
		ch.Append("contract", a.Contract)
	}

	if len(a.Action) > 0 {
		ch.Append("action", a.Action)
	}

	return ch
}
