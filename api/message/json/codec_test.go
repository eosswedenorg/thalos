package json

import (
	"testing"
	"time"

	"github.com/eosswedenorg/thalos/api/message"
	"github.com/stretchr/testify/assert"
)

func TestJson_EncodeActionTrace(t *testing.T) {
	msg := message.ActionTrace{
		TxID:      "ed3b8e853647971cf8296f004c3a1aeac255f082b2cb3c12cc3222e2d7c174ab",
		BlockNum:  267372365,
		Timestamp: time.Unix(1048267389, int64(time.Millisecond)*500).UTC(),
		Name:      "transfer",
		Contract:  "eosio",
		Receiver:  "account2",
		Data: map[string]interface{}{
			"from":     "account1",
			"to":       "account2",
			"quantity": "1000.0000 WAX",
		},
		Authorization: []message.PermissionLevel{
			{Actor: "account1", Permission: "active"},
		},
		Except: "errstr",
		Error:  2,
		Return: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	expected := `{"tx_id":"ed3b8e853647971cf8296f004c3a1aeac255f082b2cb3c12cc3222e2d7c174ab","blocknum":267372365,"blocktimestamp":"2003-03-21T17:23:09.500","name":"transfer","contract":"eosio","receiver":"account2","data":{"from":"account1","quantity":"1000.0000 WAX","to":"account2"},"authorization":[{"actor":"account1","permission":"active"}],"except":"errstr","error":2,"return":"3q2+7w=="}`

	data, err := encoder(msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(data))
}

func TestJson_DecodeActionTrace(t *testing.T) {
	expected := message.ActionTrace{
		TxID:      "952989f7464237b6cf9926e533ecd331df6794ed07564bd052bc368cbd65b4bc",
		BlockNum:  8723971,
		Timestamp: time.Unix(1718957306, int64(time.Millisecond)*500).UTC(),
		Name:      "transfer",
		Contract:  "eosio",
		Receiver:  "account2",
		Data: map[string]interface{}{
			"from":     "account1",
			"to":       "account2",
			"quantity": "1000.0000 WAX",
		},
		Authorization: []message.PermissionLevel{
			{Actor: "account1", Permission: "active"},
		},
		Except: "errstr",
		Error:  2,
		Return: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	input := `{"tx_id":"952989f7464237b6cf9926e533ecd331df6794ed07564bd052bc368cbd65b4bc","blocknum":8723971,"blocktimestamp":"2024-06-21T08:08:26.500","name":"transfer","contract":"eosio","receiver":"account2","data":{"from":"account1","quantity":"1000.0000 WAX","to":"account2"},"authorization":[{"actor":"account1","permission":"active"}],"except":"errstr","error":2,"return":"3q2+7w=="}`

	msg := message.ActionTrace{}
	err := decoder([]byte(input), &msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, msg)
}
