package api

import (
	"bytes"
	"io"
	"testing"

	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	"github.com/stretchr/testify/assert"
)

type mockReader struct {
	r io.Reader
}

func (m mockReader) Read(channel Channel) ([]byte, error) {
	if m.r != nil {
		b, err := io.ReadAll(m.r)
		if err == nil && len(b) < 1 {
			err = io.EOF
		}
		return b, err
	}
	return []byte{}, io.EOF
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

func TestClient_SubscribeWithNilHandler(t *testing.T) {
	client := NewClient(nil, nil)
	client.OnAction = mockActionHandler
	client.OnHeartbeat = mockHbHandler

	err := client.Subscribe(TableDeltaChannel{Name: "name"}.Channel())

	assert.Error(t, err)
}

func TestClient_ReadRollback(t *testing.T) {
	called := false
	expected := message.RollbackMessage{
		OldBlockNum: 1000,
		NewBlockNum: 50,
	}

	codec, err := message.GetCodec("json")
	assert.NoError(t, err)

	payload, err := codec.Encoder(expected)
	assert.NoError(t, err)

	client := NewClient(mockReader{bytes.NewReader(payload)}, codec.Decoder)
	client.OnRollback = func(rb message.RollbackMessage) {
		assert.Equal(t, rb, expected)
		called = true
	}

	err = client.Subscribe(RollbackChannel)
	assert.NoError(t, err)

	client.Run()

	assert.True(t, called, "Rollback callback not called when it should have been")
}
