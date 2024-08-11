package ship_test

import (
	"testing"

	"github.com/eosswedenorg/thalos/internal/ship"
	"github.com/stretchr/testify/assert"
)

func TestParseTableDeltaData(t *testing.T) {
	input := []interface{}{
		"resource_limits_state_v0",
		map[string]interface{}{
			"average_block_cpu_usage": []interface{}{
				"usage_accumulator_v0",
				map[string]interface{}{
					"consumed":     33679,
					"last_ordinal": 308855607,
					"value_ex":     "18525321667",
				},
			},
			"average_block_net_usage": []interface{}{
				"usage_accumulator_v0",
				map[string]interface{}{
					"consumed":     8107,
					"last_ordinal": 308855607,
					"value_ex":     int64(3854030492),
				},
			},
			"slice": []interface{}{
				"generated_transaction_v0",
				[]interface{}{1, 2, "tree"},
			},
			"single_value": []interface{}{
				"generated_transaction_v0",
				uint32(12933729),
			},
			"total_cpu_weight":  "44811223778385154",
			"total_net_weight":  "134285012330070718",
			"total_ram_bytes":   "172065109473",
			"virtual_cpu_limit": 206081,
			"virtual_net_limit": 1048576000,
		},
	}

	expected := map[string]interface{}{
		"average_block_cpu_usage": map[string]interface{}{
			"consumed":     33679,
			"last_ordinal": 308855607,
			"value_ex":     "18525321667",
		},
		"average_block_net_usage": map[string]interface{}{
			"consumed":     8107,
			"last_ordinal": 308855607,
			"value_ex":     int64(3854030492),
		},
		"slice":             []interface{}{1, 2, "tree"},
		"single_value":      uint32(12933729),
		"total_cpu_weight":  "44811223778385154",
		"total_net_weight":  "134285012330070718",
		"total_ram_bytes":   "172065109473",
		"virtual_cpu_limit": 206081,
		"virtual_net_limit": 1048576000,
	}

	actual, err := ship.ParseTableDeltaData(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
