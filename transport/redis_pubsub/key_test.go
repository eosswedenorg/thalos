package redis_pubsub

import (
	"testing"

	"eosio-ship-trace-reader/transport"
)

func TestKey_String(t *testing.T) {
	type fields struct {
		NS      Namespace
		Channel transport.ChannelInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Empty", fields{NS: Namespace{}, Channel: transport.Channel{}}, "ship::0000000000000000000000000000000000000000000000000000000000000000::"},
		{"Transactions", fields{NS: Namespace{ChainID: "id"}, Channel: transport.Channel{"transactions"}}, "ship::id::transactions"},
		{"Nested", fields{NS: Namespace{ChainID: "id"}, Channel: transport.Channel{"one.two"}}, "ship::id::one.two"},
		{"Action", fields{NS: Namespace{ChainID: "id"}, Channel: transport.ActionChannel{Contract: "mycontract"}}, "ship::id::actions/contract/mycontract"},
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
