package json

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/eosswedenorg/thalos/api/message"
	"github.com/stretchr/testify/assert"
)

func TestJson_EncodeActionTrace(t *testing.T) {
	dataJson, err := json.Marshal(map[string]string{
		"from":     "account1",
		"to":       "account2",
		"quantity": "1000.0000 WAX",
	})

	assert.NoError(t, err)

	msg := message.ActionTrace{
		TxID:      "ed3b8e853647971cf8296f004c3a1aeac255f082b2cb3c12cc3222e2d7c174ab",
		BlockNum:  267372365,
		Timestamp: time.Unix(1048267389, int64(time.Millisecond)*500).UTC(),
		Name:      "transfer",
		Contract:  "eosio",
		Receiver:  "account2",
		Data:      dataJson,
		HexData:   hex.EncodeToString(dataJson),
		Authorization: []message.PermissionLevel{
			{Actor: "account1", Permission: "active"},
		},
		Except: "errstr",
		Error:  2,
		Return: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	expected := `{"tx_id":"ed3b8e853647971cf8296f004c3a1aeac255f082b2cb3c12cc3222e2d7c174ab","blocknum":267372365,"blocktimestamp":"2003-03-21T17:23:09.500","name":"transfer","contract":"eosio","receiver":"account2","data":"eyJmcm9tIjoiYWNjb3VudDEiLCJxdWFudGl0eSI6IjEwMDAuMDAwMCBXQVgiLCJ0byI6ImFjY291bnQyIn0=","hex_data":"7b2266726f6d223a226163636f756e7431222c227175616e74697479223a22313030302e3030303020574158222c22746f223a226163636f756e7432227d","authorization":[{"actor":"account1","permission":"active"}],"except":"errstr","error":2,"return":"3q2+7w=="}`

	data, err := json_codec.Marshal(msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(data))
}

func TestJson_DecodeActionTrace(t *testing.T) {
	dataJson, err := json.Marshal(map[string]string{
		"from":     "account1",
		"to":       "account2",
		"quantity": "1000.0000 WAX",
	})

	assert.NoError(t, err)

	expected := message.ActionTrace{
		TxID:      "952989f7464237b6cf9926e533ecd331df6794ed07564bd052bc368cbd65b4bc",
		BlockNum:  8723971,
		Timestamp: time.Unix(1718957306, int64(time.Millisecond)*500).UTC(),
		Name:      "transfer",
		Contract:  "eosio",
		Receiver:  "account2",
		Data:      dataJson,
		HexData:   hex.EncodeToString(dataJson),
		Authorization: []message.PermissionLevel{
			{Actor: "account1", Permission: "active"},
		},
		Except: "errstr",
		Error:  2,
		Return: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	input := `{"tx_id":"952989f7464237b6cf9926e533ecd331df6794ed07564bd052bc368cbd65b4bc","blocknum":8723971,"blocktimestamp":"2024-06-21T08:08:26.500","name":"transfer","contract":"eosio","receiver":"account2","data":"eyJmcm9tIjoiYWNjb3VudDEiLCJxdWFudGl0eSI6IjEwMDAuMDAwMCBXQVgiLCJ0byI6ImFjY291bnQyIn0=","hex_data":"7b2266726f6d223a226163636f756e7431222c227175616e74697479223a22313030302e3030303020574158222c22746f223a226163636f756e7432227d","authorization":[{"actor":"account1","permission":"active"}],"except":"errstr","error":2,"return":"3q2+7w=="}`

	msg := message.ActionTrace{}
	err = json_codec.Unmarshal([]byte(input), &msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, msg)
}
