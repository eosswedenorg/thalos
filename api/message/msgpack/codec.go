package msgpack

import (
	"reflect"

	"github.com/ugorji/go/codec"

	"github.com/eosswedenorg/thalos/api/message"
)

func createCodec() message.Codec {
	// Create handler.
	handle := codec.MsgpackHandle{}
	handle.MapType = reflect.TypeOf(map[string]any(nil))
	handle.Canonical = true

	// Weird name but this is needed for the newest version of msgpack
	// that has support for time and string datatypes etc.
	handle.WriteExt = true

	return message.Codec{
		Encoder: func(a any) ([]byte, error) {
			var b []byte
			return b, codec.NewEncoderBytes(&b, &handle).Encode(a)
		},
		Decoder: func(b []byte, a any) error {
			return codec.NewDecoderBytes(b, &handle).Decode(a)
		},
	}
}

func init() {
	// Register codec.
	message.RegisterCodec("msgpack", createCodec())
}
