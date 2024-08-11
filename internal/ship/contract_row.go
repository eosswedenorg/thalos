package ship

import (
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
