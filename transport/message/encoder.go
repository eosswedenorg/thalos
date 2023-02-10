package message

// Encoder is a function that can encode a object to the encoded format.
type Encoder func(v any) ([]byte, error)
