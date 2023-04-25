package msgpack

import (
	"testing"

	"github.com/eosswedenorg/thalos/api/message"
	"github.com/shamaton/msgpack/v2"
	"github.com/stretchr/testify/assert"
)

func TestMsgpack_EncodeActionTrace(t *testing.T) {
	RegisterGeneratedResolver()

	msg := message.ActionTrace{
		TxID:     "edc06dce6320459fd644756972048da453b2816b0a434c37ddffde36778dcab3",
		Name:     "sellitem",
		Contract: "mygame",
		Receiver: "eosio",
		Data: map[interface{}]interface{}{
			"item": map[interface{}]interface{}{
				"id":   "2131242",
				"name": "Great Sword",
				"str":  "100",
				"agi":  "20",
				"dur":  "100",
				"qual": "epic",
			},
			"from":   "account1",
			"to":     "account2",
			"amount": "1000.0000 SCAM",
		},
		HexData: "d0fa1b2ab8a6fd0d1b0173df91aa9ffd277642d05780cf750",
	}

	data, err := msgpack.Marshal(msg)
	assert.NoError(t, err)

	res := message.ActionTrace{}
	err = msgpack.Unmarshal(data, &res)
	assert.NoError(t, err)

	assert.Equal(t, msg, res)
}

func TestMsgpack_Decode(t *testing.T) {
	RegisterGeneratedResolver()

	data := []byte("\x86\xa5tx_id\xd9@edc06dce6320459fd644756972048da453b2816b0a434c37ddffde36778dcab3\xa4name\xa8sellitem\xa8contract\xa6mygame\xa8receiver\xa5eosio\xa4data\x84\xa4item\x86\xa4name\xabGreat Sword\xa3str\xa3100\xa3agi\xa220\xa3dur\xa3100\xa4qual\xa4epic\xa2id\xa72131242\xa4from\xa8account1\xa2to\xa8account2\xa6amount\xae1000.0000 SCAM\xa8hex_data\xd91d0fa1b2ab8a6fd0d1b0173df91aa9ffd277642d05780cf750")

	expected := message.ActionTrace{
		TxID:     "edc06dce6320459fd644756972048da453b2816b0a434c37ddffde36778dcab3",
		Name:     "sellitem",
		Contract: "mygame",
		Receiver: "eosio",
		Data: map[interface{}]interface{}{
			"item": map[interface{}]interface{}{
				"id":   "2131242",
				"name": "Great Sword",
				"str":  "100",
				"agi":  "20",
				"dur":  "100",
				"qual": "epic",
			},
			"from":   "account1",
			"to":     "account2",
			"amount": "1000.0000 SCAM",
		},
		HexData: "d0fa1b2ab8a6fd0d1b0173df91aa9ffd277642d05780cf750",
	}

	res := message.ActionTrace{}
	err := msgpack.Unmarshal(data, &res)
	assert.NoError(t, err)

	assert.Equal(t, res, expected)
}

func TestMsgpack_EncodeHeartbeat(t *testing.T) {
	RegisterGeneratedResolver()

	msg := message.HeartBeat{
		BlockNum:                 1234,
		HeadBlockNum:             1235,
		LastIrreversibleBlockNum: 1236,
	}

	data, err := msgpack.Marshal(msg)
	assert.NoError(t, err)

	assert.Equal(t, data, []byte{0x83, 0xa8, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0xcd, 0x4, 0xd2, 0xad, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0xcd, 0x4, 0xd3, 0xba, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x72, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x69, 0x62, 0x6c, 0x65, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0xcd, 0x4, 0xd4})
}

func TestMsgpack_DecodeHeartbeat(t *testing.T) {
	data := []byte{0x83, 0xa8, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0xcd, 0x03, 0xe8, 0xad, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0xcd, 0x0b, 0xb8, 0xba, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69, 0x72, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x69, 0x62, 0x6c, 0x65, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0xcd, 0x04, 0x06}

	expected := message.HeartBeat{
		BlockNum:                 1000,
		HeadBlockNum:             3000,
		LastIrreversibleBlockNum: 1030,
	}

	msg := message.HeartBeat{}
	err := msgpack.Unmarshal(data, &msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, msg)
}
