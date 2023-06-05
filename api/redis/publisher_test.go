package redis

import (
	"testing"

	"github.com/eosswedenorg/thalos/api"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestPublisher_Write(t *testing.T) {
	client, mock := redismock.NewClientMock()

	pub := NewPublisher(context.Background(), client, Namespace{ChainID: "id"})

	mock.MatchExpectationsInOrder(true)
	mock.ExpectPublish("ship::id::test", []byte("some string")).SetVal(0)
	mock.ExpectPublish("ship::id::test2", []byte("some other string")).SetVal(0)

	assert.NoError(t, pub.Write(api.Channel{"test"}, []byte("some string")))
	assert.NoError(t, pub.Write(api.Channel{"test2"}, []byte("some other string")))
	assert.NoError(t, pub.Flush())

	assert.NoError(t, mock.ExpectationsWereMet())
}
