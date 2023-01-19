package redis_pubsub

import (
	"context"

	"eosio-ship-trace-reader/transport"

	redis "github.com/go-redis/redis/v8"
)

type RedisPubsub struct {
	pipeline redis.Pipeliner
	ctx      context.Context
}

func New(client *redis.Client) *RedisPubsub {
	return &RedisPubsub{
		pipeline: client.Pipeline(),
		ctx:      client.Context(),
	}
}

func (r *RedisPubsub) Publish(channel transport.ChannelInterface, payload []byte) error {
	return r.pipeline.Publish(r.ctx, channel.String(), payload).Err()
}

func (r *RedisPubsub) Flush() error {
	_, err := r.pipeline.Exec(r.ctx)
	return err
}
