package server

import (
	"fmt"
	"reflect"

	"github.com/pnx/antelope-go/ship"
)

// convert a ActionTrace to ActionTraceV1
func toActionTraceV1(trace *ship.ActionTrace) *ship.ActionTraceV1 {
	if trace.V0 != nil {
		// convert to v1
		return &ship.ActionTraceV1{
			ActionOrdinal:        trace.V0.ActionOrdinal,
			CreatorActionOrdinal: trace.V0.CreatorActionOrdinal,
			Receipt:              trace.V0.Receipt,
			Receiver:             trace.V0.Receiver,
			Act:                  trace.V0.Act,
			ContextFree:          trace.V0.ContextFree,
			Elapsed:              trace.V0.Elapsed,
			Console:              trace.V0.Console,
			AccountRamDeltas:     trace.V0.AccountRamDeltas,
			Except:               trace.V0.Except,
			ErrorCode:            trace.V0.ErrorCode,
			ReturnValue:          []byte{},
		}
	}
	return trace.V1
}

func isVariantName(name string) bool {
	validVariants := []string{
		"get_status_request_v0",
		"block_position",
		"get_status_result_v0",
		"get_blocks_request_v0",
		"get_blocks_ack_request_v0",
		"get_blocks_result_v0",
		"row",
		"table_delta_v0",
		"action",
		"account_auth_sequence",
		"action_receipt_v0",
		"account_delta",
		"action_trace_v0",
		"partial_transaction_v0",
		"transaction_trace_v0",
		"packed_transaction",
		"transaction_receipt_header",
		"transaction_receipt",
		"extension",
		"block_header",
		"signed_block_header",
		"signed_block",
		"transaction_header",
		"transaction",
		"code_id",
		"account_v0",
		"account_metadata_v0",
		"code_v0",
		"contract_table_v0",
		"contract_row_v0",
		"contract_index64_v0",
		"contract_index128_v0",
		"contract_index256_v0",
		"contract_index_double_v0",
		"contract_index_long_double_v0",
		"producer_key",
		"producer_schedule",
		"block_signing_authority_v0",
		"producer_authority",
		"producer_authority_schedule",
		"chain_config_v0",
		"global_property_v0",
		"global_property_v1",
		"generated_transaction_v0",
		"activated_protocol_feature_v0",
		"protocol_state_v0",
		"key_weight",
		"permission_level",
		"permission_level_weight",
		"wait_weight",
		"authority",
		"permission_v0",
		"permission_link_v0",
		"resource_limits_v0",
		"usage_accumulator_v0",
		"resource_usage_v0",
		"resource_limits_state_v0",
		"resource_limits_ratio_v0",
		"elastic_limit_parameters_v0",
		"resource_limits_config_v0",
	}

	for _, v := range validVariants {
		if v == name {
			return true
		}
	}
	return false
}

// Check if a structure is a variant type.
// This is not 100% accurate. As variant types comes
// as a simple slice with the types name in the first index
// and the value as the second.
// So there could be some edge cases where this structure is actual data
// and not a variant type although should be super rare.
func isVariant(v reflect.Value) bool {
	if v.Kind() != reflect.Slice || v.Len() != 2 {
		return false
	}

	for v = v.Index(0); v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer; v = v.Elem() {
	}

	return v.Kind() == reflect.String && isVariantName(v.String())
}

func parseTableDeltaData(v any) (map[string]interface{}, error) {
	iface := parseTableDeltaDataInner(reflect.ValueOf(v)).Interface()
	if out, ok := iface.(map[string]interface{}); ok {
		return out, nil
	}
	return nil, fmt.Errorf("data is not an map")
}

func parseTableDeltaDataInner(v reflect.Value) reflect.Value {
	if isVariant(v) {
		v = v.Index(1)
	}

	switch v.Kind() {
	case reflect.Interface:
		return parseTableDeltaDataInner(v.Elem())
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			v.Index(i).Set(parseTableDeltaDataInner(v.Index(i)))
		}
	case reflect.Map:
		it := v.MapRange()
		for it.Next() {
			v.SetMapIndex(it.Key(), parseTableDeltaDataInner(it.Value()))
		}
	}

	return v
}
