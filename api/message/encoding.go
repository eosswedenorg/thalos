package message

// Encoder is a function that can encode a object to the encoded format.
type Encoder func(any) ([]byte, error)

// Decoder is a function that can decode a format into an object
type Decoder func([]byte, any) error

// Codec is a type that can has a matching Encoder and Decoder function.
type Codec struct {
	Encoder Encoder
	Decoder Decoder
}
