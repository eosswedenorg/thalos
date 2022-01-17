
package main

import (
    eos "github.com/eoscanada/eos-go"
)

type ActionTrace struct {
    Receiver   eos.Name `json:"receiver"`
    Contract   eos.AccountName `json:"contract"`
	Action     eos.ActionName `json:"action"`
	Data       interface{} `json:"data"`
}
