package api

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

func TestChannel_Is(t *testing.T) {
	tests := []struct {
		name string
		a    Channel
		b    Channel
		want bool
	}{
		{"Empty valid", Channel{}, Channel{}, true},
		{"Valid #1", Channel{"one"}, Channel{"one"}, true},
		{"Valid #2", Channel{"one", "two"}, Channel{"one", "two"}, true},
		{"Invalid #1", Channel{"one"}, Channel{"one", "two"}, false},
		{"Invalid #2", Channel{"one", "three"}, Channel{"one", "two"}, false},
		{"Invalid #3", Channel{"two", "one"}, Channel{"one", "two"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Is(tt.b); got != tt.want {
				t.Errorf("a.Is(b) = %v, want %v", got, tt.want)
			}

			if got := tt.b.Is(tt.a); got != tt.want {
				t.Errorf("b.Is(a) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_Format(t *testing.T) {
	tests := []struct {
		name  string
		c     Channel
		delim string
		want  string
	}{
		{"Empty", Channel{}, ":", ""},
		{"Alot#1", Channel{"one", "two", "three"}, "-", "one-two-three"},
		{"Alot#2", Channel{"first", "second"}, ":", "first:second"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Format(tt.delim); got != tt.want {
				t.Errorf("Channel.Format() = %v, want %v", got, tt.want)
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

func TestChannel_Type(t *testing.T) {
	tests := []struct {
		name string
		c    Channel
		want string
	}{
		{"Empty", Channel{}, "unknown"},
		{"Heartbeat", HeartbeatChannel, "heartbeat"},
		{"Transaction", TransactionChannel, "transactions"},
		{"Actions", ActionChannel{}.Channel(), "actions"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Type(); got != tt.want {
				t.Errorf("Channel.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAction_Channel(t *testing.T) {
	tests := []struct {
		name   string
		action ActionChannel
		want   Channel
	}{
		{"Empty", ActionChannel{}, Channel{"actions"}},
		{"Contract", ActionChannel{Contract: "mycontract"}, Channel{"actions", "contract", "mycontract"}},
		{"Action", ActionChannel{Name: "myaction"}, Channel{"actions", "name", "myaction"}},
		{"ContractAndName", ActionChannel{Contract: "mycontract", Name: "myaction"}, Channel{"actions", "contract", "mycontract", "name", "myaction"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.Channel(); !got.Is(tt.want) {
				t.Errorf("ActionChannel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableDelta_Channel(t *testing.T) {
	tests := []struct {
		name   string
		action TableDeltaChannel
		want   Channel
	}{
		{"Empty", TableDeltaChannel{}, Channel{"tabledeltas"}},
		{"Contract", TableDeltaChannel{Name: "delta_name"}, Channel{"tabledeltas", "name", "delta_name"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.Channel(); !got.Is(tt.want) {
				t.Errorf("TableDeltaChannel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
