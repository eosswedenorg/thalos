package abi

import (
	"encoding/json"

	eos "github.com/eoscanada/eos-go"
)

func DecodeAction(eos_ABI *eos.ABI, data []byte, actionName eos.ActionName) (interface{}, error) {
	var v interface{}

	bytes, err := eos_ABI.DecodeAction(data, actionName)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(bytes, &v)
	return v, err
}
