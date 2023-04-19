package redis_common

import "testing"

func TestNamespace_String(t *testing.T) {
	tests := []struct {
		name string
		ns   Namespace
		want string
	}{
		{"Empty", Namespace{}, "ship::0000000000000000000000000000000000000000000000000000000000000000"},
		{"Prefix Only", Namespace{Prefix: "some.prefix"}, "some.prefix::0000000000000000000000000000000000000000000000000000000000000000"},
		{"ChainID Only", Namespace{ChainID: "1234"}, "ship::1234"},
		{"Both", Namespace{Prefix: "my.prefix", ChainID: "1234"}, "my.prefix::1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ns.String(); got != tt.want {
				t.Errorf("Namespace.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
