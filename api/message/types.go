package message

import (
	"errors"
	"time"
)

type HeartBeat struct {
	BlockNum                 uint32 `json:"blocknum" msgpack:"blocknum"`
	HeadBlockNum             uint32 `json:"head_blocknum" msgpack:"head_blocknum"`
	LastIrreversibleBlockNum uint32 `json:"last_irreversible_blocknum" msgpack:"last_irreversible_blocknum"`
}

type PermissionLevel struct {
	Actor      string `json:"actor" msgpack:"actor"`
	Permission string `json:"permission" msgpack:"permission"`
}

type AccountAuthSequence struct {
	Account  string `json:"account" msgpack:"account"`
	Sequence uint64 `json:"sequence" msgpack:"sequence"`
}

type TransactionTrace struct {
	ID            string        `json:"id" msgpack:"id"`
	BlockNum      uint32        `json:"blocknum" msgpack:"blocknum"`
	Timestamp     time.Time     `json:"blocktimestamp" msgpack:"blocktimestamp"`
	Status        string        `json:"status" msgpack:"status"`
	CPUUsageUS    uint32        `json:"cpu_usage_us" msgpack:"cpu_usage_us"`
	NetUsageWords uint32        `json:"net_usage_words" msgpack:"net_usage_words"`
	Elapsed       int64         `json:"elapsed" msgpack:"elapsed"`
	NetUsage      uint64        `json:"net_usage" msgpack:"net_usage"`
	Scheduled     bool          `json:"scheduled" msgpack:"scheduled"`
	ActionTraces  []ActionTrace `json:"action_traces" msgpack:"action_traces"`
	// AccountDelta    *eos.AccountRAMDelta `json:"account_delta" eos:"optional"`
	Except string `json:"except" msgpack:"except"`
	Error  uint64 `json:"error" msgpack:"error"`
	// FailedDtrxTrace *TransactionTrace   `json:"failed_dtrx_trace" eos:"optional"`
	// Partial         *PartialTransaction `json:"partial" eos:"optional"`
}

type ActionReceipt struct {
	Receiver       string                `json:"receiver" msgpack:"receiver"`
	ActDigest      string                `json:"act_digest" msgpack:"act_digest"`
	GlobalSequence uint64                `json:"global_sequence" msgpack:"global_sequence"`
	RecvSequence   uint64                `json:"recv_sequence" msgpack:"recv_sequence"`
	AuthSequence   []AccountAuthSequence `json:"auth_sequence" msgpack:"auth_sequence"`
	CodeSequence   uint32                `json:"code_sequence" msgpack:"code_sequence"`
	ABISequence    uint32                `json:"abi_sequence" msgpack:"abi_sequence"`
}

type ActionTrace struct {
	TxID string `json:"tx_id" msgpack:"tx_id"`

	BlockNum uint32 `json:"blocknum" msgpack:"blocknum"`

	Timestamp time.Time `json:"blocktimestamp" msgpack:"blocktimestamp"`

	Receipt *ActionReceipt `json:"receipt,omitempty" msgpack:"receipt"`

	// Action name
	Name string `json:"name" msgpack:"name"`

	// Contract account.
	Contract string `json:"contract" msgpack:"contract"`

	Receiver string      `json:"receiver" msgpack:"receiver"`
	Data     interface{} `json:"data" msgpack:"data"`

	Authorization []PermissionLevel `json:"authorization" msgpack:"authorization"`

	Except string `json:"except" msgpack:"except"`
	Error  uint64 `json:"error" msgpack:"error"`
	Return []byte `json:"return" msgpack:"return"`
}

func (act ActionTrace) GetData() (map[string]any, error) {
	if data, ok := act.Data.(map[string]any); ok {
		return data, nil
	}
	return nil, errors.New("failed to convert data to map")
}

type TableDeltaRow struct {
	Present bool   `json:"present" msgpack:"present"`
	Data    []byte `json:"data" msgpack:"data"`
	RawData []byte `json:"raw_data" msgpack:"raw_data"`
}

type TableDelta struct {
	BlockNum  uint32          `json:"blocknum" msgpack:"blocknum"`
	Timestamp time.Time       `json:"blocktimestamp" msgpack:"blocktimestamp"`
	Name      string          `json:"name" msgpack:"name"`
	Rows      []TableDeltaRow `json:"rows" msgpack:"rows"`
}
