package redis_common

import (
	"testing"

	"thalos/api"
)

func TestKey_String(t *testing.T) {
	type fields struct {
		NS      Namespace
		Channel api.Channel
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Empty", fields{NS: Namespace{}, Channel: api.Channel{}}, "ship::0000000000000000000000000000000000000000000000000000000000000000::"},
		{"Transactions", fields{NS: Namespace{ChainID: "id"}, Channel: api.Channel{"transactions"}}, "ship::id::transactions"},
		{"Nested", fields{NS: Namespace{ChainID: "id"}, Channel: api.Channel{"one.two"}}, "ship::id::one.two"},
		{"Action", fields{NS: Namespace{ChainID: "id"}, Channel: api.Action{Contract: "mycontract"}.Channel()}, "ship::id::actions/contract/mycontract"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := Key{
				NS:      tt.fields.NS,
				Channel: tt.fields.Channel,
			}
			if got := k.String(); got != tt.want {
				t.Errorf("Key.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
