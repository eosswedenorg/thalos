package json

import (
	"time"

	"github.com/eosswedenorg/thalos/api/message"
	"github.com/shufflingpixels/jsontime-go"
)

func createCodec() message.Codec {
	json_codec := jsontime.New("2006-01-02T15:04:05.000", time.UTC)

	return message.Codec{
		Encoder: json_codec.Marshal,
		Decoder: json_codec.Unmarshal,
	}
}

func init() {
	// Register the json codec.
	message.RegisterCodec("json", createCodec())
}
