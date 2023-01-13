package redis_pubsub

import (
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisPubsub(t *testing.T) {
	client, mock := redismock.NewClientMock()

	pubsub := New(client)

	mock.MatchExpectationsInOrder(true)
	mock.ExpectPublish("test", []byte("some string")).SetVal(0)
	mock.ExpectPublish("test2", []byte("some other string")).SetVal(0)

	assert.NoError(t, pubsub.Publish("test", []byte("some string")))
	assert.NoError(t, pubsub.Publish("test2", []byte("some other string")))
	assert.NoError(t, pubsub.Flush())

	assert.NoError(t, mock.ExpectationsWereMet())
}
