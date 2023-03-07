package redis_pubsub

import (
	"context"

	"eosio-ship-trace-reader/transport"

	"github.com/go-redis/redis/v8"
)

type Publisher struct {
	pipeline redis.Pipeliner
	ctx      context.Context
	ns       Namespace
}

func NewPublisher(client *redis.Client, ns Namespace) *Publisher {
	return &Publisher{
		pipeline: client.Pipeline(),
		ctx:      client.Context(),
		ns:       ns,
	}
}

func (r *Publisher) Write(channel transport.ChannelInterface, payload []byte) error {
	return r.pipeline.Publish(r.ctx, r.ns.NewKey(channel).String(), payload).Err()
}

func (r *Publisher) Flush() error {
	_, err := r.pipeline.Exec(r.ctx)
	return err
}
