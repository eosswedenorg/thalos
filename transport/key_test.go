package transport

import "testing"

func TestKey_String(t *testing.T) {
	type fields struct {
		NS      Namespace
		Channel ChannelInterface
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Empty", fields{NS: Namespace{}, Channel: Channel{}}, "ship::0000000000000000000000000000000000000000000000000000000000000000::"},
		{"Transactions", fields{NS: Namespace{ChainID: "id"}, Channel: Channel{"transactions"}}, "ship::id::transactions"},
		{"Nested", fields{NS: Namespace{ChainID: "id"}, Channel: Channel{"one.two"}}, "ship::id::one.two"},
		{"Action", fields{NS: Namespace{ChainID: "id"}, Channel: ActionChannel{Contract: "mycontract"}}, "ship::id::actions/contract/mycontract"},
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
