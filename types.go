
package main

import (
    eos "github.com/eoscanada/eos-go"
)

type ActionTrace struct {
    Receiver   eos.Name
    Contract   eos.AccountName
	Action     eos.ActionName
	Data       interface{}
}
