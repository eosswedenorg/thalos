package main

import (
	eos "github.com/eoscanada/eos-go"
)

type ActionTrace struct {
	TxID     eos.Checksum256 `json:"tx_id"`
	Receiver eos.Name        `json:"receiver"`
	Contract eos.AccountName `json:"contract"`
	Action   eos.ActionName  `json:"action"`
	Data     interface{}     `json:"data"`
	HexData  string          `json:"hex_data"`
}
