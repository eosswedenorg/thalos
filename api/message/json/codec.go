package json

import (
	"time"

	jsontime "github.com/eosswedenorg-go/jsontime/v2"
	"github.com/eosswedenorg/thalos/api/message"
)

var json_codec = jsontime.ConfigWithCustomTimeFormat

func init() {
	jsontime.SetDefaultTimeFormat("2006-01-02T15:04:05.000", time.UTC)

	message.RegisterCodec("json", message.Codec{
		Encoder: json_codec.Marshal,
		Decoder: json_codec.Unmarshal,
	})
}
