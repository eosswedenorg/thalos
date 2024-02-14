package redis

import (
	"context"

	"github.com/eosswedenorg/thalos/api"
	. "github.com/eosswedenorg/thalos/api/redis"

	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	pipeline redis.Pipeliner
	ctx      context.Context
	ns       Namespace
}

func NewPublisher(ctx context.Context, client *redis.Client, ns Namespace) *Publisher {
	return &Publisher{
		pipeline: client.Pipeline(),
		ctx:      ctx,
		ns:       ns,
	}
}

func (r *Publisher) Write(channel api.Channel, payload []byte) error {
	return r.pipeline.Publish(r.ctx, r.ns.NewKey(channel).String(), payload).Err()
}

func (r *Publisher) Flush() error {
	_, err := r.pipeline.Exec(r.ctx)
	return err
}

func (r *Publisher) Close() error {
	return r.Flush()
}
