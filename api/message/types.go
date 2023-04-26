package message

import "encoding/json"

type HeartBeat struct {
	BlockNum                 uint32 `json:"blocknum" msgpack:"blocknum"`
	HeadBlockNum             uint32 `json:"head_blocknum" msgpack:"head_blocknum"`
	LastIrreversibleBlockNum uint32 `json:"last_irreversible_blocknum" msgpack:"last_irreversible_blocknum"`
}

type ActionTrace struct {
	TxID string `json:"tx_id" msgpack:"tx_id"`

	// Action name
	Name string `json:"name" msgpack:"name"`

	// Contract account.
	Contract string `json:"contract" msgpack:"contract"`

	Receiver string `json:"receiver" msgpack:"receiver"`
	Data     []byte `json:"data" msgpack:"data"`
	HexData  string `json:"hex_data" msgpack:"hex_data"`
}

func (act ActionTrace) GetData() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	return data, json.Unmarshal(act.Data, &data)
}
