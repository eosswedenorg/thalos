package api

// Reader interface defines the required method
// to read a message from an channel.
//
// This is a low-level interface typically implemented by backend drivers
type Reader interface {
	// Read a message from a channel.
	// Read may block until a message is ready or an error occurred.
	//
	// io.EOF is returned from a reader when there is no more data to be read.
	// If Read returns io.EOF all subsequent calls must also return io.EOF
	//
	// This function should be designed to handle concurrent calls. eg. thread safe.
	Read(channel Channel) ([]byte, error)

	// Close closes the reader
	// Any blocked Read operations will be unblocked and return io.EOF
	Close() error
}
