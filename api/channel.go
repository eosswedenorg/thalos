package api

import (
	"strings"
)

// Channel is just a wrapper around string slice
type Channel []string

func (c *Channel) Append(name ...string) {
	*c = append(*c, name...)
}

func (c Channel) Format(delimiter string) string {
	return strings.Join(c, delimiter)
}

func (c Channel) String() string {
	return c.Format("/")
}

func (c Channel) Type() string {
	if len(c) > 0 {
		return c[0]
	}
	return "unknown"
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
	TransactionChannel = Channel{"transactions"}
	HeartbeatChannel   = Channel{"heartbeat"}
)

// Action Channel
type ActionChannel struct {
	Name     string
	Contract string
}

func (a ActionChannel) Channel() Channel {
	ch := Channel{"actions"}

	if len(a.Contract) > 0 {
		ch.Append("contract", a.Contract)
	}

	if len(a.Name) > 0 {
		ch.Append("name", a.Name)
	}

	return ch
}

// Table deltas
type TableDeltaChannel struct {
	Name string
}

func (td TableDeltaChannel) Channel() Channel {
	ch := Channel{"tabledeltas"}

	if len(td.Name) > 0 {
		ch.Append("name", td.Name)
	}

	return ch
}
