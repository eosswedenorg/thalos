package ship

import "reflect"

func IsVariantName(name string) bool {
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
func IsVariant(v reflect.Value) bool {
	if v.Kind() != reflect.Slice || v.Len() != 2 {
		return false
	}

	for v = v.Index(0); v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer; v = v.Elem() {
		// Intentionally empty
	}

	return v.Kind() == reflect.String && IsVariantName(v.String())
}
