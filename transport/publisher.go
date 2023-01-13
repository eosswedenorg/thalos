package transport

// Publisher interface defines the required methods
// to send messages over channel.
type Publisher interface {
	// Publish a message to a channel.
	// The message may or may not be buffered depending on the implementation.
	Publish(channel string, payload []byte) error

	// Flush writes any buffered messages to the channel.
	// If the implementation does not support buffering. this is a noop.
	Flush() error
}
