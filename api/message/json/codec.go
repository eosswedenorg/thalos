package json

import (
	jsontime "github.com/eosswedenorg-go/jsontime/v2"
	"github.com/eosswedenorg/thalos/api/message"
)

func createCodec() message.Codec {
	json_codec := jsontime.ConfigWithCustomTimeFormat

	return message.Codec{
		Encoder: json_codec.Marshal,
		Decoder: json_codec.Unmarshal,
	}
}

func init() {
	// Register the json codec.
	message.RegisterCodec("json", createCodec())
}
