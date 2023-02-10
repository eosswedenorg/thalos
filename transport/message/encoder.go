package message

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

// Encoder is a function that can encode a object to the encoded format.
type Encoder func(v any) ([]byte, error)

func Encode(v interface{}) ([]byte, bool) {
	payload, err := json.Marshal(v)
	if err != nil {
		log.WithError(err).
			WithField("v", v).
			Warn("Failed to encode message to json")
		return nil, false
	}
	return payload, true
}
