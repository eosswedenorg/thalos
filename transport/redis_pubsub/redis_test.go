package redis_pubsub

import (
	"testing"

	"eosio-ship-trace-reader/transport"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisPubsub(t *testing.T) {
	client, mock := redismock.NewClientMock()

	pubsub := New(client, Namespace{ChainID: "id"})

	mock.MatchExpectationsInOrder(true)
	mock.ExpectPublish("ship::id::test", []byte("some string")).SetVal(0)
	mock.ExpectPublish("ship::id::test2", []byte("some other string")).SetVal(0)

	assert.NoError(t, pubsub.Publish(transport.Channel{"test"}, []byte("some string")))
	assert.NoError(t, pubsub.Publish(transport.Channel{"test2"}, []byte("some other string")))
	assert.NoError(t, pubsub.Flush())

	assert.NoError(t, mock.ExpectationsWereMet())
}
