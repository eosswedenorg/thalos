package msgpack

import (
	"testing"
	"time"

	"github.com/eosswedenorg/thalos/api/message"
	"github.com/stretchr/testify/assert"
)

func TestMsgpack_EncodeActionTrace(t *testing.T) {
	msg := message.ActionTrace{
		TxID:      "edc06dce6320459fd644756972048da453b2816b0a434c37ddffde36778dcab3",
		BlockNum:  12345,
		Timestamp: time.Unix(1699617279, int64(time.Millisecond)*500),
		Name:      "sellitem",
		Contract:  "mygame",
		Receiver:  "eosio",
		Receipt: &message.ActionReceipt{
			Receiver:       "eosio",
			ActDigest:      "be9618d12f0b8d125731c6faf1304357291ada716bb190c6c03c8dd41c36bb79",
			GlobalSequence: 273871,
			RecvSequence:   237863,
			AuthSequence: []message.AccountAuthSequence{
				{
					Account:  "eosio",
					Sequence: 217328173,
				},
			},
			CodeSequence: 2381267931,
			ABISequence:  847623,
		},
		Data: map[string]any{
			"item": map[string]any{
				"id":   2131242,
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
		Authorization: []message.PermissionLevel{
			{Actor: "mygame", Permission: "active"},
		},
		Except: "errstr",
		Error:  2,
		Return: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	data, err := createCodec().Encoder(msg)
	assert.NoError(t, err)

	expected := []byte{
		0x8d, 0xad, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
		0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x91,
		0x82, 0xa5, 0x61, 0x63, 0x74, 0x6f, 0x72, 0xa6,
		0x6d, 0x79, 0x67, 0x61, 0x6d, 0x65, 0xaa, 0x70,
		0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
		0x6e, 0xa6, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65,
		0xa8, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75,
		0x6d, 0xcd, 0x30, 0x39, 0xae, 0x62, 0x6c, 0x6f,
		0x63, 0x6b, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
		0x61, 0x6d, 0x70, 0xd7, 0xff, 0x77, 0x35, 0x94,
		0x00, 0x65, 0x4e, 0x19, 0xff, 0xa8, 0x63, 0x6f,
		0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0xa6, 0x6d,
		0x79, 0x67, 0x61, 0x6d, 0x65, 0xa4, 0x64, 0x61,
		0x74, 0x61, 0x84, 0xa6, 0x61, 0x6d, 0x6f, 0x75,
		0x6e, 0x74, 0xae, 0x31, 0x30, 0x30, 0x30, 0x2e,
		0x30, 0x30, 0x30, 0x30, 0x20, 0x53, 0x43, 0x41,
		0x4d, 0xa4, 0x66, 0x72, 0x6f, 0x6d, 0xa8, 0x61,
		0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x31, 0xa4,
		0x69, 0x74, 0x65, 0x6d, 0x86, 0xa3, 0x61, 0x67,
		0x69, 0xa2, 0x32, 0x30, 0xa3, 0x64, 0x75, 0x72,
		0xa3, 0x31, 0x30, 0x30, 0xa2, 0x69, 0x64, 0xd2,
		0x00, 0x20, 0x85, 0x2a, 0xa4, 0x6e, 0x61, 0x6d,
		0x65, 0xab, 0x47, 0x72, 0x65, 0x61, 0x74, 0x20,
		0x53, 0x77, 0x6f, 0x72, 0x64, 0xa4, 0x71, 0x75,
		0x61, 0x6c, 0xa4, 0x65, 0x70, 0x69, 0x63, 0xa3,
		0x73, 0x74, 0x72, 0xa3, 0x31, 0x30, 0x30, 0xa2,
		0x74, 0x6f, 0xa8, 0x61, 0x63, 0x63, 0x6f, 0x75,
		0x6e, 0x74, 0x32, 0xa5, 0x65, 0x72, 0x72, 0x6f,
		0x72, 0x02, 0xa6, 0x65, 0x78, 0x63, 0x65, 0x70,
		0x74, 0xa6, 0x65, 0x72, 0x72, 0x73, 0x74, 0x72,
		0xae, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x72,
		0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0xc2,
		0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa8, 0x73, 0x65,
		0x6c, 0x6c, 0x69, 0x74, 0x65, 0x6d, 0xa7, 0x72,
		0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x87, 0xac,
		0x61, 0x62, 0x69, 0x5f, 0x73, 0x65, 0x71, 0x75,
		0x65, 0x6e, 0x63, 0x65, 0xce, 0x00, 0x0c, 0xef,
		0x07, 0xaa, 0x61, 0x63, 0x74, 0x5f, 0x64, 0x69,
		0x67, 0x65, 0x73, 0x74, 0xd9, 0x40, 0x62, 0x65,
		0x39, 0x36, 0x31, 0x38, 0x64, 0x31, 0x32, 0x66,
		0x30, 0x62, 0x38, 0x64, 0x31, 0x32, 0x35, 0x37,
		0x33, 0x31, 0x63, 0x36, 0x66, 0x61, 0x66, 0x31,
		0x33, 0x30, 0x34, 0x33, 0x35, 0x37, 0x32, 0x39,
		0x31, 0x61, 0x64, 0x61, 0x37, 0x31, 0x36, 0x62,
		0x62, 0x31, 0x39, 0x30, 0x63, 0x36, 0x63, 0x30,
		0x33, 0x63, 0x38, 0x64, 0x64, 0x34, 0x31, 0x63,
		0x33, 0x36, 0x62, 0x62, 0x37, 0x39, 0xad, 0x61,
		0x75, 0x74, 0x68, 0x5f, 0x73, 0x65, 0x71, 0x75,
		0x65, 0x6e, 0x63, 0x65, 0x91, 0x82, 0xa7, 0x61,
		0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0xa5, 0x65,
		0x6f, 0x73, 0x69, 0x6f, 0xa8, 0x73, 0x65, 0x71,
		0x75, 0x65, 0x6e, 0x63, 0x65, 0xce, 0x0c, 0xf4,
		0x2a, 0x2d, 0xad, 0x63, 0x6f, 0x64, 0x65, 0x5f,
		0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65,
		0xce, 0x8d, 0xef, 0x43, 0xdb, 0xaf, 0x67, 0x6c,
		0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x73, 0x65, 0x71,
		0x75, 0x65, 0x6e, 0x63, 0x65, 0xce, 0x00, 0x04,
		0x2d, 0xcf, 0xa8, 0x72, 0x65, 0x63, 0x65, 0x69,
		0x76, 0x65, 0x72, 0xa5, 0x65, 0x6f, 0x73, 0x69,
		0x6f, 0xad, 0x72, 0x65, 0x63, 0x76, 0x5f, 0x73,
		0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0xce,
		0x00, 0x03, 0xa1, 0x27, 0xa8, 0x72, 0x65, 0x63,
		0x65, 0x69, 0x76, 0x65, 0x72, 0xa5, 0x65, 0x6f,
		0x73, 0x69, 0x6f, 0xa6, 0x72, 0x65, 0x74, 0x75,
		0x72, 0x6e, 0xc4, 0x04, 0xde, 0xad, 0xbe, 0xef,
		0xa5, 0x74, 0x78, 0x5f, 0x69, 0x64, 0xd9, 0x40,
		0x65, 0x64, 0x63, 0x30, 0x36, 0x64, 0x63, 0x65,
		0x36, 0x33, 0x32, 0x30, 0x34, 0x35, 0x39, 0x66,
		0x64, 0x36, 0x34, 0x34, 0x37, 0x35, 0x36, 0x39,
		0x37, 0x32, 0x30, 0x34, 0x38, 0x64, 0x61, 0x34,
		0x35, 0x33, 0x62, 0x32, 0x38, 0x31, 0x36, 0x62,
		0x30, 0x61, 0x34, 0x33, 0x34, 0x63, 0x33, 0x37,
		0x64, 0x64, 0x66, 0x66, 0x64, 0x65, 0x33, 0x36,
		0x37, 0x37, 0x38, 0x64, 0x63, 0x61, 0x62, 0x33,
	}

	assert.Equal(t, expected, data)
}

func TestMsgpack_DecodeActionTrace(t *testing.T) {
	data := []byte("\x8c\xadauthorization\x91\x82\xa5actor\xa6mygame\xaapermission\xa6active\xa8blocknum\xce\x00\x85F7\xaeblocktimestamp\xd6\xffH\xf1U\x1f\xa8contract\xa6mygame\xa4data\x83\xafdropped_from_id\xd2\x00\nK\x02\xa4item\x86\xa3dur\xd1\x00\x91\xa2id\xd2\x00\x00\xc1פname\xacShadowmourne\xa4qual\xa9legendary\xa3sta\xd1\x00ƣstr\xd1\x00ߨreceiver\xa8account1\xa5error\x02\xa6except\xa6errstr\xa4name\xa4drop\xa7receipt\x87\xacabi_sequence\xce\x00\bӽ\xaaact_digest\xd9@676c5336de9d528d456b01e975b1006e2bdc86c8d566330321e9b309b634f523\xadauth_sequence\x91\x82\xa7account\xa6mygame\xa8sequence\xce\x00\v\v\x10\xadcode_sequence\xce\x0e2\x82\xb3\xafglobal_sequence\xce\x03r`6\xa8receiver\xa8account1\xadrecv_sequence\xce\x00\x13\x81S\xa8receiver\xa8account1\xa6return\xc4\x04ޭ\xbe\xef\xa5tx_id\xd9@edc06dce6320459fd644756972048da453b2816b0a434c37ddffde36778dcab3")

	expected := message.ActionTrace{
		TxID:      "edc06dce6320459fd644756972048da453b2816b0a434c37ddffde36778dcab3",
		BlockNum:  8734263,
		Timestamp: time.Unix(1223775519, 0).UTC(),
		Name:      "drop",
		Contract:  "mygame",
		Receiver:  "account1",
		Receipt: &message.ActionReceipt{
			Receiver:       "account1",
			ActDigest:      "676c5336de9d528d456b01e975b1006e2bdc86c8d566330321e9b309b634f523",
			GlobalSequence: 57827382,
			RecvSequence:   1278291,
			AuthSequence: []message.AccountAuthSequence{
				{
					Account:  "mygame",
					Sequence: 723728,
				},
			},
			CodeSequence: 238191283,
			ABISequence:  578493,
		},
		Data: map[string]any{
			"item": map[string]any{
				"id":   int64(49623),
				"name": "Shadowmourne",
				"str":  int64(223),
				"sta":  int64(198),
				"dur":  int64(145),
				"qual": "legendary",
			},
			"dropped_from_id": int64(674562),
			"receiver":        "account1",
		},
		Authorization: []message.PermissionLevel{
			{Actor: "mygame", Permission: "active"},
		},
		Except: "errstr",
		Error:  2,
		Return: []byte{0xde, 0xad, 0xbe, 0xef},
	}

	res := message.ActionTrace{}
	err := createCodec().Decoder(data, &res)
	assert.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestMsgpack_EncodeHeartbeat(t *testing.T) {
	msg := message.HeartBeat{
		BlockNum:                 1234,
		HeadBlockNum:             1235,
		LastIrreversibleBlockNum: 1236,
	}

	data, err := createCodec().Encoder(msg)
	assert.NoError(t, err)

	assert.Equal(t, data, []byte{
		0x83, 0xa8, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e,
		0x75, 0x6d, 0xcd, 0x04, 0xd2, 0xad, 0x68, 0x65,
		0x61, 0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
		0x6e, 0x75, 0x6d, 0xcd, 0x04, 0xd3, 0xba, 0x6c,
		0x61, 0x73, 0x74, 0x5f, 0x69, 0x72, 0x72, 0x65,
		0x76, 0x65, 0x72, 0x73, 0x69, 0x62, 0x6c, 0x65,
		0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75,
		0x6d, 0xcd, 0x04, 0xd4,
	})
}

func TestMsgpack_DecodeHeartbeat(t *testing.T) {
	data := []byte{
		0x83, 0xa8, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e,
		0x75, 0x6d, 0xcd, 0x03, 0xe8, 0xad, 0x68, 0x65,
		0x61, 0x64, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
		0x6e, 0x75, 0x6d, 0xcd, 0x0b, 0xb8, 0xba, 0x6c,
		0x61, 0x73, 0x74, 0x5f, 0x69, 0x72, 0x72, 0x65,
		0x76, 0x65, 0x72, 0x73, 0x69, 0x62, 0x6c, 0x65,
		0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75,
		0x6d, 0xcd, 0x04, 0x06,
	}

	expected := message.HeartBeat{
		BlockNum:                 1000,
		HeadBlockNum:             3000,
		LastIrreversibleBlockNum: 1030,
	}

	msg := message.HeartBeat{}
	err := createCodec().Decoder(data, &msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, msg)
}

func TestMsgpack_EncodeTransactionTrace(t *testing.T) {
	tx_id := "edc06dce6320459fd644756972048da453b2816b0a434c37ddffde36778dcab3"
	block_num := uint32(2738723781)
	ts := time.Unix(1699617279, int64(time.Millisecond)*500).UTC()

	msg := message.TransactionTrace{
		ID:            tx_id,
		BlockNum:      block_num,
		Timestamp:     ts,
		Status:        "soft_fail",
		CPUUsageUS:    23,
		NetUsageWords: 16,
		NetUsage:      128,
		Elapsed:       4,
		Scheduled:     true,
		ActionTraces: []message.ActionTrace{
			{
				Name:     "sellitem",
				Contract: "skjdh23",
				Receiver: "eosio",
				Receipt: &message.ActionReceipt{
					Receiver:       "eosio",
					ActDigest:      "676c5336de9d528d456b01e975b1006e2bdc86c8d566330321e9b309b634f523",
					GlobalSequence: 34789213,
					RecvSequence:   578912378,
					AuthSequence: []message.AccountAuthSequence{
						{
							Account:  "actor12123",
							Sequence: 547897,
						},
					},
					CodeSequence: 23193721,
					ABISequence:  4782232,
				},
				Data: map[string]any{
					"key": "value",
				},
				Authorization: []message.PermissionLevel{
					{Actor: "actor12123", Permission: "active"},
				},
			},
			{
				Name:     "sellitem",
				Contract: "skjdh24",
				Receiver: "eosio",
				Authorization: []message.PermissionLevel{
					{Actor: "actor1482", Permission: "active"},
				},
			},
		},
		Except: "exceptstr",
		Error:  2,
	}

	data, err := createCodec().Encoder(msg)
	assert.NoError(t, err)

	expected := []byte{
		0x8c, 0xad, 0x61, 0x63, 0x74, 0x69, 0x6f,
		0x6e, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65,
		0x73, 0x92, 0x8d, 0xad, 0x61, 0x75, 0x74,
		0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
		0x69, 0x6f, 0x6e, 0x91, 0x82, 0xa5, 0x61,
		0x63, 0x74, 0x6f, 0x72, 0xaa, 0x61, 0x63,
		0x74, 0x6f, 0x72, 0x31, 0x32, 0x31, 0x32,
		0x33, 0xaa, 0x70, 0x65, 0x72, 0x6d, 0x69,
		0x73, 0x73, 0x69, 0x6f, 0x6e, 0xa6, 0x61,
		0x63, 0x74, 0x69, 0x76, 0x65, 0xa8, 0x62,
		0x6c, 0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d,
		0x00, 0xae, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
		0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
		0x6d, 0x70, 0xc0, 0xa8, 0x63, 0x6f, 0x6e,
		0x74, 0x72, 0x61, 0x63, 0x74, 0xa7, 0x73,
		0x6b, 0x6a, 0x64, 0x68, 0x32, 0x33, 0xa4,
		0x64, 0x61, 0x74, 0x61, 0x81, 0xa3, 0x6b,
		0x65, 0x79, 0xa5, 0x76, 0x61, 0x6c, 0x75,
		0x65, 0xa5, 0x65, 0x72, 0x72, 0x6f, 0x72,
		0x00, 0xa6, 0x65, 0x78, 0x63, 0x65, 0x70,
		0x74, 0xa0, 0xae, 0x66, 0x69, 0x72, 0x73,
		0x74, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69,
		0x76, 0x65, 0x72, 0xc2, 0xa4, 0x6e, 0x61,
		0x6d, 0x65, 0xa8, 0x73, 0x65, 0x6c, 0x6c,
		0x69, 0x74, 0x65, 0x6d, 0xa7, 0x72, 0x65,
		0x63, 0x65, 0x69, 0x70, 0x74, 0x87, 0xac,
		0x61, 0x62, 0x69, 0x5f, 0x73, 0x65, 0x71,
		0x75, 0x65, 0x6e, 0x63, 0x65, 0xce, 0x00,
		0x48, 0xf8, 0x98, 0xaa, 0x61, 0x63, 0x74,
		0x5f, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74,
		0xd9, 0x40, 0x36, 0x37, 0x36, 0x63, 0x35,
		0x33, 0x33, 0x36, 0x64, 0x65, 0x39, 0x64,
		0x35, 0x32, 0x38, 0x64, 0x34, 0x35, 0x36,
		0x62, 0x30, 0x31, 0x65, 0x39, 0x37, 0x35,
		0x62, 0x31, 0x30, 0x30, 0x36, 0x65, 0x32,
		0x62, 0x64, 0x63, 0x38, 0x36, 0x63, 0x38,
		0x64, 0x35, 0x36, 0x36, 0x33, 0x33, 0x30,
		0x33, 0x32, 0x31, 0x65, 0x39, 0x62, 0x33,
		0x30, 0x39, 0x62, 0x36, 0x33, 0x34, 0x66,
		0x35, 0x32, 0x33, 0xad, 0x61, 0x75, 0x74,
		0x68, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65,
		0x6e, 0x63, 0x65, 0x91, 0x82, 0xa7, 0x61,
		0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0xaa,
		0x61, 0x63, 0x74, 0x6f, 0x72, 0x31, 0x32,
		0x31, 0x32, 0x33, 0xa8, 0x73, 0x65, 0x71,
		0x75, 0x65, 0x6e, 0x63, 0x65, 0xce, 0x00,
		0x08, 0x5c, 0x39, 0xad, 0x63, 0x6f, 0x64,
		0x65, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65,
		0x6e, 0x63, 0x65, 0xce, 0x01, 0x61, 0xe8,
		0x79, 0xaf, 0x67, 0x6c, 0x6f, 0x62, 0x61,
		0x6c, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65,
		0x6e, 0x63, 0x65, 0xce, 0x02, 0x12, 0xd7,
		0x5d, 0xa8, 0x72, 0x65, 0x63, 0x65, 0x69,
		0x76, 0x65, 0x72, 0xa5, 0x65, 0x6f, 0x73,
		0x69, 0x6f, 0xad, 0x72, 0x65, 0x63, 0x76,
		0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e,
		0x63, 0x65, 0xce, 0x22, 0x81, 0x80, 0x7a,
		0xa8, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
		0x65, 0x72, 0xa5, 0x65, 0x6f, 0x73, 0x69,
		0x6f, 0xa6, 0x72, 0x65, 0x74, 0x75, 0x72,
		0x6e, 0xc0, 0xa5, 0x74, 0x78, 0x5f, 0x69,
		0x64, 0xa0, 0x8c, 0xad, 0x61, 0x75, 0x74,
		0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
		0x69, 0x6f, 0x6e, 0x91, 0x82, 0xa5, 0x61,
		0x63, 0x74, 0x6f, 0x72, 0xa9, 0x61, 0x63,
		0x74, 0x6f, 0x72, 0x31, 0x34, 0x38, 0x32,
		0xaa, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73,
		0x73, 0x69, 0x6f, 0x6e, 0xa6, 0x61, 0x63,
		0x74, 0x69, 0x76, 0x65, 0xa8, 0x62, 0x6c,
		0x6f, 0x63, 0x6b, 0x6e, 0x75, 0x6d, 0x00,
		0xae, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x74,
		0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
		0x70, 0xc0, 0xa8, 0x63, 0x6f, 0x6e, 0x74,
		0x72, 0x61, 0x63, 0x74, 0xa7, 0x73, 0x6b,
		0x6a, 0x64, 0x68, 0x32, 0x34, 0xa4, 0x64,
		0x61, 0x74, 0x61, 0xc0, 0xa5, 0x65, 0x72,
		0x72, 0x6f, 0x72, 0x00, 0xa6, 0x65, 0x78,
		0x63, 0x65, 0x70, 0x74, 0xa0, 0xae, 0x66,
		0x69, 0x72, 0x73, 0x74, 0x5f, 0x72, 0x65,
		0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0xc2,
		0xa4, 0x6e, 0x61, 0x6d, 0x65, 0xa8, 0x73,
		0x65, 0x6c, 0x6c, 0x69, 0x74, 0x65, 0x6d,
		0xa8, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76,
		0x65, 0x72, 0xa5, 0x65, 0x6f, 0x73, 0x69,
		0x6f, 0xa6, 0x72, 0x65, 0x74, 0x75, 0x72,
		0x6e, 0xc0, 0xa5, 0x74, 0x78, 0x5f, 0x69,
		0x64, 0xa0, 0xa8, 0x62, 0x6c, 0x6f, 0x63,
		0x6b, 0x6e, 0x75, 0x6d, 0xce, 0xa3, 0x3d,
		0x9b, 0xc5, 0xae, 0x62, 0x6c, 0x6f, 0x63,
		0x6b, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
		0x61, 0x6d, 0x70, 0xd7, 0xff, 0x77, 0x35,
		0x94, 0x00, 0x65, 0x4e, 0x19, 0xff, 0xac,
		0x63, 0x70, 0x75, 0x5f, 0x75, 0x73, 0x61,
		0x67, 0x65, 0x5f, 0x75, 0x73, 0x17, 0xa7,
		0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64,
		0x04, 0xa5, 0x65, 0x72, 0x72, 0x6f, 0x72,
		0x02, 0xa6, 0x65, 0x78, 0x63, 0x65, 0x70,
		0x74, 0xa9, 0x65, 0x78, 0x63, 0x65, 0x70,
		0x74, 0x73, 0x74, 0x72, 0xa2, 0x69, 0x64,
		0xd9, 0x40, 0x65, 0x64, 0x63, 0x30, 0x36,
		0x64, 0x63, 0x65, 0x36, 0x33, 0x32, 0x30,
		0x34, 0x35, 0x39, 0x66, 0x64, 0x36, 0x34,
		0x34, 0x37, 0x35, 0x36, 0x39, 0x37, 0x32,
		0x30, 0x34, 0x38, 0x64, 0x61, 0x34, 0x35,
		0x33, 0x62, 0x32, 0x38, 0x31, 0x36, 0x62,
		0x30, 0x61, 0x34, 0x33, 0x34, 0x63, 0x33,
		0x37, 0x64, 0x64, 0x66, 0x66, 0x64, 0x65,
		0x33, 0x36, 0x37, 0x37, 0x38, 0x64, 0x63,
		0x61, 0x62, 0x33, 0xa9, 0x6e, 0x65, 0x74,
		0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0xcc,
		0x80, 0xaf, 0x6e, 0x65, 0x74, 0x5f, 0x75,
		0x73, 0x61, 0x67, 0x65, 0x5f, 0x77, 0x6f,
		0x72, 0x64, 0x73, 0x10, 0xa9, 0x73, 0x63,
		0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64,
		0xc3, 0xa6, 0x73, 0x74, 0x61, 0x74, 0x75,
		0x73, 0xa9, 0x73, 0x6f, 0x66, 0x74, 0x5f,
		0x66, 0x61, 0x69, 0x6c,
	}
	assert.Equal(t, expected, data)
}

func TestMsgpack_DecodeTransactionTrace(t *testing.T) {
	data := []byte("\x8c\xadaction_traces\x91\x8b\xadauthorization\x91\x82\xa5actor\xaaactor12123\xaapermission\xa5claim\xa8blocknum\x00\xaeblocktimestamp\xc0\xa8contract\xa9send.help\xa4data\x81\xa7claimer\xaaactor12123\xa5error\x00\xa6except\xa0\xa4name\xa5claim\xa8receiver\xaaactor12123\xa6return\xc0\xa5tx_id\xa0\xa8blocknumΣ5\x81!\xaeblocktimestamp\xd7\xffw5\x94\x00K\x813̬cpu_usage_us6\xa7elapsed\x04\xa5error+\xa6except\xa9exceptstr\xa2id\xd9@05d7e50e8aa898a84df345f714f741ce804a9cc171da44b893ae74891cc7258a\xa9net_usagè\xafnet_usage_words\x10\xa9scheduledæstatus\xa9hard_fail")

	tx_id := "05d7e50e8aa898a84df345f714f741ce804a9cc171da44b893ae74891cc7258a"
	block_num := uint32(2738192673)
	ts := time.Unix(1266758604, int64(time.Millisecond)*500).UTC()

	expected := message.TransactionTrace{
		ID:            tx_id,
		BlockNum:      block_num,
		Timestamp:     ts,
		Status:        "hard_fail",
		CPUUsageUS:    54,
		NetUsageWords: 16,
		NetUsage:      128,
		Elapsed:       4,
		Scheduled:     true,
		ActionTraces: []message.ActionTrace{
			{
				Name:     "claim",
				Contract: "send.help",
				Receiver: "actor12123",
				Data: map[string]any{
					"claimer": "actor12123",
				},
				Authorization: []message.PermissionLevel{
					{Actor: "actor12123", Permission: "claim"},
				},
			},
		},
		Except: "exceptstr",
		Error:  43,
	}

	res := message.TransactionTrace{}
	err := createCodec().Decoder(data, &res)
	assert.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestMsgpack_EncodeTableDelta(t *testing.T) {
	msg := message.TableDelta{
		BlockNum:  6347293,
		Timestamp: time.Date(1998, time.December, 4, 8, 54, 35, 500, time.UTC),
		Name:      "contract_row",
		Rows: []message.TableDeltaRow{
			{
				Present: true,
				Data: map[string]any{
					"id":   2213,
					"name": "Freddie Mercury",
					"band": "Queen",
				},
				RawData: []byte{0x23, 0x13, 0xe2},
			},
			{
				Present: false,
				Data: map[string]any{
					"id":   27182,
					"name": "Eddie Van Halen",
					"band": "Van Halen",
				},
				RawData: []byte{0xfe, 0x4e, 0x52, 0x05},
			},
		},
	}

	expected := "\x84\xa8blocknum\xce\x00`\xda\x1d\xaeblocktimestamp\xd7\xff\x00\x00\a\xd06g\xa3K\xa4name\xaccontract_row\xa4rows\x92\x83\xa4data\x83\xa4band\xa5Queen\xa2id\xd1\b\xa5\xa4name\xafFreddie Mercury\xa7presentèraw_data\xc4\x03#\x13⃤data\x83\xa4band\xa9Van Halen\xa2id\xd1j.\xa4name\xafEddie Van Halen\xa7present¨raw_data\xc4\x04\xfeNR\x05"

	actual, err := createCodec().Encoder(msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestMsgpack_DecodeTableDelta(t *testing.T) {
	data := []byte("\x84\xa8blocknum\xce\x00`\xda\x1d\xaeblocktimestamp\xd7\xff\x00\x00\a\xd06g\xa3K\xa4name\xaccontract_row\xa4rows\x92\x83\xa4data\x83\xa4band\xa5Queen\xa2id\xd1\b\xa5\xa4name\xafFreddie Mercury\xa7presentèraw_data\xc4\x03#\x13⃤data\x83\xa4band\xa9Van Halen\xa2id\xd1j.\xa4name\xafEddie Van Halen\xa7present¨raw_data\xc4\x04\xfeNR\x05")

	expected := message.TableDelta{
		BlockNum:  6347293,
		Timestamp: time.Date(1998, time.December, 4, 8, 54, 35, 500, time.UTC),
		Name:      "contract_row",
		Rows: []message.TableDeltaRow{
			{
				Present: true,
				Data: map[string]any{
					"id":   int64(2213),
					"name": "Freddie Mercury",
					"band": "Queen",
				},
				RawData: []byte{0x23, 0x13, 0xe2},
			},
			{
				Present: false,
				Data: map[string]any{
					"id":   int64(27182),
					"name": "Eddie Van Halen",
					"band": "Van Halen",
				},
				RawData: []byte{0xfe, 0x4e, 0x52, 0x05},
			},
		},
	}
	actual := message.TableDelta{}
	err := createCodec().Decoder(data, &actual)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
