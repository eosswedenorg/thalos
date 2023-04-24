package json

import (
	"encoding/json"

	"github.com/eosswedenorg/thalos/api/message"
)

func init() {
	message.RegisterCodec("json", message.Codec{
		Encoder: json.Marshal,
		Decoder: json.Unmarshal,
	})
}
