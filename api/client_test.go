package api

import (
	"testing"

	"github.com/eosswedenorg/thalos/api/message"
)

type mockReader struct{}

func (m mockReader) Read(channel Channel) ([]byte, error) {
	return []byte{}, nil
}

func (m mockReader) Close() error {
	return nil
}

func mockDecoder([]byte, any) error {
	return nil
}

func mockHbHandler(message.HeartBeat) {
}

func mockActionHandler(message.ActionTrace) {
}

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
			c := NewClient(&mockReader{}, mockDecoder)
			c.OnHeartbeat = mockHbHandler
			c.OnAction = mockActionHandler
			if err := c.Subscribe(tt.channel); (err != nil) != tt.wantErr {
				t.Errorf("Client.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
