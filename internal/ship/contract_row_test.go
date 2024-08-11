package ship_test

import (
	"testing"

	"github.com/eosswedenorg/thalos/internal/ship"
	"github.com/shufflingpixels/antelope-go/chain"
	"github.com/stretchr/testify/assert"
)

func TestContractRow_Decode(t *testing.T) {
	expected := &ship.ContractRow{
		Code:       chain.N("eosio"),
		Scope:      chain.N("scope"),
		Table:      chain.N("accounts"),
		PrimaryKey: "1278127812",
		Payer:      chain.N("account1"),
		Value:      []byte{0x01, 0x01, 0x02, 0x03},
	}

	actual, err := ship.DecodeContractRow(map[string]any{
		"code":        uint64(6138663577826885632),
		"scope":       uint64(13990807175891517440),
		"table":       uint64(3607749779137757184),
		"primary_key": uint32(1278127812),
		"payer":       uint64(3607749778751881216),
		"value":       []byte{0x01, 0x01, 0x02, 0x03},
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
