package ship

import (
	"bytes"

	"github.com/eosswedenorg/thalos/internal/abi"
	"github.com/mitchellh/mapstructure"
	"github.com/shufflingpixels/antelope-go/chain"
)

type ContractRow struct {
	Code       chain.Name  `mapstructure:"code"`
	Scope      chain.Name  `mapstructure:"scope"`
	Table      chain.Name  `mapstructure:"table"`
	PrimaryKey string      `mapstructure:"primary_key"`
	Payer      chain.Name  `mapstructure:"payer"`
	Value      chain.Bytes `mapstructure:"value"`
}

func ParseContractRow(v map[string]interface{}) (*ContractRow, error) {
	out := &ContractRow{}
	err := mapstructure.WeakDecode(v, out)
	return out, err
}

func DecodeContractRow(manager *abi.AbiManager, data map[string]any) (any, error) {
	row, err := ParseContractRow(data)
	if err != nil {
		return nil, err
	}

	abi, err := manager.GetAbi(row.Code)
	if err != nil {
		return nil, err
	}
	return abi.DecodeTable(bytes.NewReader(row.Value), row.Table)
}
