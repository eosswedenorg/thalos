package msgpack

import (
	"github.com/shamaton/msgpack/v2"

	"github.com/eosswedenorg/thalos/api/message"
)

//go:generate go run github.com/shamaton/msgpackgen -v -input-file ../types.go -output-file msgpack.go

func init() {
	RegisterGeneratedResolver()

	message.RegisterCodec("msgpack", message.Codec{
		Encoder: msgpack.Marshal,
		Decoder: msgpack.Unmarshal,
	})
}
