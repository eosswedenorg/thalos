package redis_pubsub

import (
	"testing"

	"thalos/transport"
	. "thalos/transport/redis_common"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_Write(t *testing.T) {
	client, mock := redismock.NewClientMock()

	pub := NewPublisher(client, Namespace{ChainID: "id"})

	mock.MatchExpectationsInOrder(true)
	mock.ExpectPublish("ship::id::test", []byte("some string")).SetVal(0)
	mock.ExpectPublish("ship::id::test2", []byte("some other string")).SetVal(0)

	assert.NoError(t, pub.Write(transport.Channel{"test"}, []byte("some string")))
	assert.NoError(t, pub.Write(transport.Channel{"test2"}, []byte("some other string")))
	assert.NoError(t, pub.Flush())

	assert.NoError(t, mock.ExpectationsWereMet())
}
