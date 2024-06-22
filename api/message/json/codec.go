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
	// Set timeformat used by SHIP api
	jsontime.AddTimeFormatAlias("ship", "2006-01-02T15:04:05.000")
	jsontime.AddLocaleAlias("ship", time.UTC)

	// Register the json codec.
	message.RegisterCodec("json", createCodec())
}
