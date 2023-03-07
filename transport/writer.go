package transport

// Writer interface defines the required methods
// to send messages over an channel.
type Writer interface {
	// Write writes a message over a channel.
	// The message may or may not be buffered depending on the implementation.
	Write(channel Channel, payload []byte) error

	// Flush writes any buffered messages to the channel.
	// If the implementation does not support buffering. this is a noop.
	Flush() error
}
