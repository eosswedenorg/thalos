package redis

import (
	"testing"
	"time"

	"github.com/eosswedenorg/thalos/api"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestSubscriber_Construct(t *testing.T) {
	client, _ := redismock.NewClientMock()
	ns := Namespace{Prefix: "prefix", ChainID: "8f2f6ec19400d372c9b3340b1438e9c805cf9e69be962fa81d055bc037ceed8d"}

	s := NewSubscriber(client, ns)

	assert.Equal(t, s.client, client)
	assert.Equal(t, s.ctx, client.Context())
	assert.NotNil(t, s.sub)
	assert.Equal(t, s.ns, ns)
	assert.Equal(t, s.timeout, 200*time.Millisecond)

	s = NewSubscriber(client, ns, WithTimeout(4*time.Second))
	assert.Equal(t, s.timeout, 4*time.Second)
}

func TestSubscriber_Read(t *testing.T) {
	expectedMessages := []string{"payload", "payload2", "payload3"}

	server := miniredis.RunT(t)

	client := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	s := NewSubscriber(client, Namespace{Prefix: "prefix", ChainID: "d41dbd2921d5a377325661427090c6c508904d60920d6b7ea771c58da5299754"})

	go func() {
		time.Sleep(time.Millisecond * 10)

		for _, msg := range expectedMessages {
			server.Publish("prefix::d41dbd2921d5a377325661427090c6c508904d60920d6b7ea771c58da5299754::test", msg)
		}
	}()

	// Redis pubsub does not guarentee that messages are sent in the correct order.
	for range expectedMessages {
		msg, err := s.Read(api.Channel{"test"})
		assert.NoError(t, err)

		assert.Contains(t, expectedMessages, string(msg))
	}
}
