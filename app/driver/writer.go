package driver

import "github.com/eosswedenorg/thalos/api"

// Writer interface defines the required methods
// to send messages over an channel.
//
// This is a low-level interface typically implemented by backend drivers
type Writer interface {
	// Write writes a message over a channel.
	// The message may or may not be buffered depending on the implementation.
	Write(channel api.Channel, payload []byte) error

	// Flush writes any buffered messages to the channel.
	// If the implementation does not support buffering. this is a noop.
	Flush() error

	// Close closes the writer
	// Any blocked Flush or Write operations will be unblocked.
	Close() error
}
