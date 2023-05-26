package api

import (
	"testing"
)

func TestClient_Subscribe(t *testing.T) {
	tests := []struct {
		name    string
		channel Channel
		wantErr bool
	}{
		{"Channel", Channel{}, true},
		{"ActionChannel", ActionChannel{}.Channel(), false},
		{"HeartbeatChannel", HeartbeatChannel, false},
		{"TransactionChannel", TransactionChannel, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{}
			if err := c.Subscribe(tt.channel); (err != nil) != tt.wantErr {
				t.Errorf("Client.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
