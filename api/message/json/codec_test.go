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
		Receipt: &message.ActionReceipt{
			Receiver:       "account2",
			ActDigest:      "0a2c71dba327cf13a107d3a4f91c9c98f510a8fbbb483b233e222033f13a3e36",
			GlobalSequence: 2329381932,
			RecvSequence:   22,
			AuthSequence: []message.AccountAuthSequence{
				{
					Account:  "account1",
					Sequence: 1382772,
				},
			},
			CodeSequence: 1122,
			ABISequence:  12352,
		},
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

	expected := `{"tx_id":"ed3b8e853647971cf8296f004c3a1aeac255f082b2cb3c12cc3222e2d7c174ab","blocknum":267372365,"blocktimestamp":"2003-03-21T17:23:09.500","receipt":{"receiver":"account2","act_digest":"0a2c71dba327cf13a107d3a4f91c9c98f510a8fbbb483b233e222033f13a3e36","global_sequence":2329381932,"recv_sequence":22,"auth_sequence":[{"account":"account1","sequence":1382772}],"code_sequence":1122,"abi_sequence":12352},"name":"transfer","contract":"eosio","receiver":"account2","data":{"from":"account1","quantity":"1000.0000 WAX","to":"account2"},"authorization":[{"actor":"account1","permission":"active"}],"except":"errstr","error":2,"return":"3q2+7w=="}`

	data, err := createCodec().Encoder(msg)
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
		Receipt: &message.ActionReceipt{
			Receiver:       "account2",
			ActDigest:      "f2f682847fd5bf00fb315b075dc00b4ff0ce18776758077b86a233dea49dc047",
			GlobalSequence: 287328,
			RecvSequence:   42,
			AuthSequence: []message.AccountAuthSequence{
				{
					Account:  "account1",
					Sequence: 271877283,
				},
			},
			CodeSequence: 237823,
			ABISequence:  68323,
		},
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

	input := `{"tx_id":"952989f7464237b6cf9926e533ecd331df6794ed07564bd052bc368cbd65b4bc","blocknum":8723971,"blocktimestamp":"2024-06-21T08:08:26.500","receipt":{"receiver":"account2","act_digest":"f2f682847fd5bf00fb315b075dc00b4ff0ce18776758077b86a233dea49dc047","global_sequence":287328,"recv_sequence":42,"auth_sequence":[{"account":"account1","sequence":271877283}],"code_sequence":237823,"abi_sequence":68323},"name":"transfer","contract":"eosio","receiver":"account2","data":{"from":"account1","quantity":"1000.0000 WAX","to":"account2"},"authorization":[{"actor":"account1","permission":"active"}],"except":"errstr","error":2,"return":"3q2+7w=="}`

	msg := message.ActionTrace{}
	err := createCodec().Decoder([]byte(input), &msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, msg)
}

func TestJson_EncodeTransactionTrace(t *testing.T) {
	tx_id := "ed04516bdd1194aa5f0ab4c8c5445eec542c17f45a85bb3e9e4bc33e1a2486f8"
	ts := time.Unix(1865257805, int64(time.Millisecond)*500).UTC()
	block_num := uint32(283781923)

	msg := message.TransactionTrace{
		Timestamp:     ts,
		BlockNum:      block_num,
		ID:            tx_id,
		Status:        "executed",
		CPUUsageUS:    442,
		NetUsage:      128,
		NetUsageWords: 16,
		Elapsed:       22,
		Scheduled:     true,
		ActionTraces: []message.ActionTrace{
			{
				TxID:      tx_id,
				BlockNum:  block_num,
				Timestamp: ts,
				Receiver:  "actor01",
				Contract:  "coolgame",
				Name:      "mine",
				Authorization: []message.PermissionLevel{
					{
						Actor:      "actor01",
						Permission: "active",
					},
				},
				Data: map[string]any{
					"equipment_id": 1234,
					"location_id":  5445453,
				},
				Return: []byte{0x08, 0xf1},
			},
			{
				TxID:      tx_id,
				BlockNum:  block_num,
				Timestamp: ts,
				Receiver:  "coolgame",
				Contract:  "coolgame",
				Name:      "addpoints",
				Authorization: []message.PermissionLevel{
					{
						Actor:      "coolgame",
						Permission: "usrpoints",
					},
				},
				Data: map[string]any{
					"points": "1023.0423 SCAM",
				},
				Error:  2,
				Except: "some error string",
				Return: []byte{0xff, 0x02},
			},
		},
		Except: "errstr",
		Error:  2,
	}

	expected := `{"id":"ed04516bdd1194aa5f0ab4c8c5445eec542c17f45a85bb3e9e4bc33e1a2486f8","blocknum":283781923,"blocktimestamp":"2029-02-08T15:10:05.500","status":"executed","cpu_usage_us":442,"net_usage_words":16,"elapsed":22,"net_usage":128,"scheduled":true,"action_traces":[{"tx_id":"ed04516bdd1194aa5f0ab4c8c5445eec542c17f45a85bb3e9e4bc33e1a2486f8","blocknum":283781923,"blocktimestamp":"2029-02-08T15:10:05.500","name":"mine","contract":"coolgame","receiver":"actor01","data":{"equipment_id":1234,"location_id":5445453},"authorization":[{"actor":"actor01","permission":"active"}],"except":"","error":0,"return":"CPE="},{"tx_id":"ed04516bdd1194aa5f0ab4c8c5445eec542c17f45a85bb3e9e4bc33e1a2486f8","blocknum":283781923,"blocktimestamp":"2029-02-08T15:10:05.500","name":"addpoints","contract":"coolgame","receiver":"coolgame","data":{"points":"1023.0423 SCAM"},"authorization":[{"actor":"coolgame","permission":"usrpoints"}],"except":"some error string","error":2,"return":"/wI="}],"except":"errstr","error":2}`

	data, err := createCodec().Encoder(msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(data))
}

func TestJson_DecodeTransactionTrace(t *testing.T) {
	tx_id := "f58bf8a0137fcea644dbc2b0cc5b6a017a848cd33b2e924703e7e3c6d1ca0c2e"
	ts := time.Unix(1730755743, int64(time.Millisecond)*500).UTC()
	block_num := uint32(2378197231)

	input := `{"id":"f58bf8a0137fcea644dbc2b0cc5b6a017a848cd33b2e924703e7e3c6d1ca0c2e","blocknum":2378197231,"blocktimestamp":"2024-11-04T21:29:03.500","status":"executed","cpu_usage_us":442,"net_usage_words":16,"elapsed":22,"net_usage":128,"scheduled":true,"action_traces":[{"tx_id":"f58bf8a0137fcea644dbc2b0cc5b6a017a848cd33b2e924703e7e3c6d1ca0c2e","blocknum":2378197231,"blocktimestamp":"2024-11-04T21:29:03.500","receipt":{"receiver":"actor01","act_digest":"bf68a227a617aa9138215d31a2847005175f01151c6feeaeb6513fe2d090ff3a","global_sequence":47322,"recv_sequence":23289,"auth_sequence":[{"account":"actor01","sequence":56637232}],"code_sequence":98272,"abi_sequence":6823762},"name":"mine","contract":"","receiver":"actor01","data":{"equipment_id":1234,"location_id":5445453},"authorization":[{"actor":"actor01","permission":"active"}],"except":"","error":2,"return":"AQI="},{"tx_id":"f58bf8a0137fcea644dbc2b0cc5b6a017a848cd33b2e924703e7e3c6d1ca0c2e","blocknum":2378197231,"blocktimestamp":"2024-11-04T21:29:03.500","name":"addpoints","contract":"","receiver":"coolgame","data":{"points":"1023.0423 SCAM"},"authorization":[{"actor":"coolgame","permission":"usrpoints"}],"except":"","error":2,"return":"CPE="}],"except":"errstr","error":2}`

	expected := message.TransactionTrace{
		Timestamp:     ts,
		BlockNum:      block_num,
		ID:            tx_id,
		Status:        "executed",
		CPUUsageUS:    442,
		NetUsage:      128,
		NetUsageWords: 16,
		Elapsed:       22,
		Scheduled:     true,
		ActionTraces: []message.ActionTrace{
			{
				TxID:      tx_id,
				BlockNum:  block_num,
				Timestamp: ts,
				Receiver:  "actor01",
				Name:      "mine",
				Receipt: &message.ActionReceipt{
					Receiver:       "actor01",
					ActDigest:      "bf68a227a617aa9138215d31a2847005175f01151c6feeaeb6513fe2d090ff3a",
					GlobalSequence: 47322,
					RecvSequence:   23289,
					AuthSequence: []message.AccountAuthSequence{
						{
							Account:  "actor01",
							Sequence: 56637232,
						},
					},
					CodeSequence: 98272,
					ABISequence:  6823762,
				},
				Authorization: []message.PermissionLevel{
					{
						Actor:      "actor01",
						Permission: "active",
					},
				},
				Data: map[string]any{
					"equipment_id": float64(1234),
					"location_id":  float64(5445453),
				},
				Error:  2,
				Return: []byte{0x01, 0x02},
			},
			{
				TxID:      tx_id,
				BlockNum:  block_num,
				Timestamp: ts,
				Receiver:  "coolgame",
				Name:      "addpoints",
				Authorization: []message.PermissionLevel{
					{
						Actor:      "coolgame",
						Permission: "usrpoints",
					},
				},
				Data: map[string]any{
					"points": "1023.0423 SCAM",
				},
				Error:  2,
				Return: []byte{0x08, 0xf1},
			},
		},
		Except: "errstr",
		Error:  2,
	}

	msg := message.TransactionTrace{}
	err := createCodec().Decoder([]byte(input), &msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, msg)
}

func TestJson_EncodeTableDelta(t *testing.T) {
	msg := message.TableDelta{
		BlockNum:  2736712863,
		Timestamp: time.Unix(1724791632, 0).UTC(),
		Name:      "contract_row",
		Rows: []message.TableDeltaRow{
			{
				Present: true,
				Data: map[string]any{
					"string": "value",
					"int":    1234,
				},
				RawData: []byte{0x23, 0x54, 0xF2, 0xab, 0xeb},
			},
			{
				Present: false,
				Data: map[string]any{
					"string": "value",
					"int":    1234,
				},
				RawData: []byte{0x23, 0xe2, 0x3c, 0xb9},
			},
		},
	}

	expected := `{"blocknum":2736712863,"blocktimestamp":"2024-08-27T20:47:12.000","name":"contract_row","rows":[{"present":true,"data":{"int":1234,"string":"value"},"raw_data":"I1Tyq+s="},{"present":false,"data":{"int":1234,"string":"value"},"raw_data":"I+I8uQ=="}]}`

	data, err := createCodec().Encoder(msg)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(data))
}

func TestJson_DecodeTableDelta(t *testing.T) {
	data := `{"blocknum":2787282732,"blocktimestamp":"2014-02-18T04:20:43.500","name":"contract_row","rows":[{"present":true,"data":{"id":2,"name":"some_name"},"raw_data":"XbosHA1WDbI="},{"present":false,"data":{"id":1234,"name":"some_other_name"},"raw_data":"KNQA1g=="}]}`

	expected := message.TableDelta{
		BlockNum:  2787282732,
		Timestamp: time.Date(2014, time.February, 18, 4, 20, 43, 500000000, time.UTC),
		Name:      "contract_row",
		Rows: []message.TableDeltaRow{
			{
				Present: true,
				Data: map[string]any{
					"name": "some_name",
					"id":   float64(2),
				},
				RawData: []byte{0x5d, 0xba, 0x2c, 0x1c, 0x0d, 0x56, 0x0d, 0xb2},
			},
			{
				Present: false,
				Data: map[string]any{
					"name": "some_other_name",
					"id":   float64(1234),
				},
				RawData: []byte{0x28, 0xd4, 0x00, 0xd6},
			},
		},
	}

	actual := message.TableDelta{}
	err := createCodec().Decoder([]byte(data), &actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
