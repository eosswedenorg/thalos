package json

import (
	"time"

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
	// Set timeformat used by eos api.
	jsontime.SetDefaultTimeFormat("2006-01-02T15:04:05.000", time.UTC)

	// Register the json codec.
	message.RegisterCodec("json", createCodec())
}
