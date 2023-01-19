package app

import (
	eos "github.com/eoscanada/eos-go"
)

type HearthBeat struct {
	BlockNum                 uint32 `json:"blocknum"`
	HeadBlockNum             uint32 `json:"head_blocknum"`
	LastIrreversibleBlockNum uint32 `json:"last_irreversible_blocknum"`
}

type ActionTrace struct {
	TxID     eos.Checksum256 `json:"tx_id"`
	Receiver eos.Name        `json:"receiver"`
	Contract eos.AccountName `json:"contract"`
	Action   eos.ActionName  `json:"action"`
	Data     interface{}     `json:"data"`
	HexData  string          `json:"hex_data"`
}
