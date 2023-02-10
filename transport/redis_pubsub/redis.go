package redis_pubsub

import (
	"context"

	"eosio-ship-trace-reader/transport"

	"github.com/go-redis/redis/v8"
)

type RedisPubsub struct {
	pipeline redis.Pipeliner
	ctx      context.Context
	ns       Namespace
}

func New(client *redis.Client, ns Namespace) *RedisPubsub {
	return &RedisPubsub{
		pipeline: client.Pipeline(),
		ctx:      client.Context(),
		ns:       ns,
	}
}

func (r *RedisPubsub) Publish(channel transport.ChannelInterface, payload []byte) error {
	return r.pipeline.Publish(r.ctx, r.ns.NewKey(channel).String(), payload).Err()
}

func (r *RedisPubsub) Flush() error {
	_, err := r.pipeline.Exec(r.ctx)
	return err
}
