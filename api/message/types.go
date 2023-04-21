package message

type HeartBeat struct {
	BlockNum                 uint32 `json:"blocknum"`
	HeadBlockNum             uint32 `json:"head_blocknum"`
	LastIrreversibleBlockNum uint32 `json:"last_irreversible_blocknum"`
}

type ActionTrace struct {
	TxID string `json:"tx_id"`

	// Action name
	Name string `json:"name"`

	// Contract account.
	Contract string `json:"contract"`

	Receiver string      `json:"receiver"`
	Data     interface{} `json:"data"`
	HexData  string      `json:"hex_data"`
}
