package redis_stream

import (
	"errors"
	"testing"

	"eosio-ship-trace-reader/transport"
	. "eosio-ship-trace-reader/transport/redis_common"

	_redis "github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestWriter_Construct(t *testing.T) {
	client, _ := redismock.NewClientMock()

	w := NewWriter(client)
	assert.Equal(t, w.max_len, int64(2000))
	assert.Equal(t, w.ns, Namespace{})

	w = NewWriter(client, WithNamespace(Namespace{ChainID: "4422"}))
	assert.Equal(t, w.max_len, int64(2000))
	assert.Equal(t, w.ns.ChainID, "4422")

	w = NewWriter(client, WithNamespace(Namespace{ChainID: "id"}), WithMaxLen(4000))
	assert.Equal(t, w.max_len, int64(4000))
	assert.Equal(t, w.ns.ChainID, "id")
}

func TestWriter_Write(t *testing.T) {
	client, mock := redismock.NewClientMock()

	w := NewWriter(client, WithNamespace(Namespace{ChainID: "id"}))

	mock.MatchExpectationsInOrder(true)
	mock.ExpectXAdd(&_redis.XAddArgs{Stream: "ship::id::test", MaxLen: 2000, Values: []byte("some string")}).SetVal("OK")
	mock.ExpectXAdd(&_redis.XAddArgs{Stream: "ship::id::test2", MaxLen: 2000, Values: []byte("some other string")}).SetVal("OK")

	assert.NoError(t, w.Write(transport.Channel{"test"}, []byte("some string")))
	assert.NoError(t, w.Write(transport.Channel{"test2"}, []byte("some other string")))
	assert.NoError(t, w.Flush())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWriter_Write_Error(t *testing.T) {
	client, mock := redismock.NewClientMock()

	w := NewWriter(client, WithNamespace(Namespace{ChainID: "1234"}))

	mock.MatchExpectationsInOrder(true)
	mock.ExpectXAdd(&_redis.XAddArgs{Stream: "ship::1234::test", MaxLen: 2000, Values: []byte("message")}).SetErr(errors.New("ErrTestValue"))

	assert.NoError(t, w.Write(transport.Channel{"test"}, []byte("message")))
	assert.EqualError(t, w.Flush(), "ErrTestValue")

	assert.NoError(t, mock.ExpectationsWereMet())
}
