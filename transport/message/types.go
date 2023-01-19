package message

type HearthBeat struct {
	BlockNum                 uint32 `json:"blocknum"`
	HeadBlockNum             uint32 `json:"head_blocknum"`
	LastIrreversibleBlockNum uint32 `json:"last_irreversible_blocknum"`
}

type ActionTrace struct {
	TxID     string      `json:"tx_id"`
	Receiver string      `json:"receiver"`
	Contract string      `json:"contract"`
	Action   string      `json:"action"`
	Data     interface{} `json:"data"`
	HexData  string      `json:"hex_data"`
}
