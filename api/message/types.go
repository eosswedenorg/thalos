package message

import (
	"errors"
	"time"
)

type HeartBeat struct {
	BlockNum                 uint32 `json:"blocknum" msgpack:"blocknum"`
	HeadBlockNum             uint32 `json:"head_blocknum" msgpack:"head_blocknum"`
	LastIrreversibleBlockNum uint32 `json:"last_irreversible_blocknum" msgpack:"last_irreversible_blocknum"`
}

type PermissionLevel struct {
	Actor      string `json:"actor" msgpack:"actor"`
	Permission string `json:"permission" msgpack:"permission"`
}

type ActionTrace struct {
	TxID string `json:"tx_id" msgpack:"tx_id"`

	BlockNum uint32 `json:"blocknum" msgpack:"blocknum"`

	Timestamp time.Time `json:"blocktimestamp" msgpack:"blocktimestamp"`

	// Action name
	Name string `json:"name" msgpack:"name"`

	// Contract account.
	Contract string `json:"contract" msgpack:"contract"`

	Receiver string      `json:"receiver" msgpack:"receiver"`
	Data     interface{} `json:"data" msgpack:"data"`

	Authorization []PermissionLevel `json:"authorization" msgpack:"authorization"`

	Except string `json:"except" msgpack:"except"`
	Error  uint64 `json:"error" msgpack:"error"`
	Return []byte `json:"return" msgpack:"return"`
}

func (act ActionTrace) GetData() (map[string]any, error) {
	if data, ok := act.Data.(map[string]any); ok {
		return data, nil
	}
	return nil, errors.New("failed to convert data to map")
}
