package transport

// Reader interface defines the required method
// to read a message from an channel.
//
// This is a low-level interface typically implemented by transport drivers.
type Reader interface {
	// Read a message from a channel.
	// Read may block until a message is ready or an error occured.
	Read(channel ChannelInterface) ([]byte, error)
}
