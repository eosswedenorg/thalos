package transport

import (
	"reflect"
	"testing"
)

func TestChannel_Append(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		obj      Channel
		expected Channel
	}{
		{"One", "one", Channel{}, Channel{"one"}},
		{"More", "more", Channel{"one", "two"}, Channel{"one", "two", "more"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.obj.Append(tt.arg)
			if reflect.DeepEqual(tt.obj, tt.expected) == false {
				t.Errorf("Channel.Append() expected %v, got %v", tt.expected, tt.obj)
			}
		})
	}
}

func TestChannel_String(t *testing.T) {
	tests := []struct {
		name string
		c    Channel
		want string
	}{
		{"Empty", Channel{}, ""},
		{"Alot", Channel{"one", "two", "three"}, "one/two/three"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Channel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionChannel_String(t *testing.T) {
	type fields struct {
		Contract string
		Action   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Empty", fields{}, "actions"},
		{"Contract", fields{Contract: "mycontract"}, "actions/contract/mycontract"},
		{"Action", fields{Action: "myaction"}, "actions/action/myaction"},
		{"ContractAction", fields{Contract: "mycontract", Action: "myaction"}, "actions/contract/mycontract/action/myaction"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := ActionChannel{
				Contract: tt.fields.Contract,
				Action:   tt.fields.Action,
			}
			if got := ac.String(); got != tt.want {
				t.Errorf("ActionChannel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
