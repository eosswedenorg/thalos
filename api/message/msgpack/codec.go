package msgpack

import (
	"reflect"

	"github.com/ugorji/go/codec"

	"github.com/eosswedenorg/thalos/api/message"
)

var mh codec.MsgpackHandle

func encode(a any) ([]byte, error) {
	var b []byte
	return b, codec.NewEncoderBytes(&b, &mh).Encode(a)
}

func decode(b []byte, a any) error {
	return codec.NewDecoderBytes(b, &mh).Decode(a)
}

func init() {
	mh.MapType = reflect.TypeOf(map[string]any(nil))
	mh.Canonical = true

	// Wierd name but this is needed for the newest version of msgpack
	// that has support for time and string datatypes etc.
	mh.WriteExt = true

	message.RegisterCodec("msgpack", message.Codec{
		Encoder: encode,
		Decoder: decode,
	})
}
