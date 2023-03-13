package redis_stream

import (
	"context"

	"eosio-ship-trace-reader/transport"
	. "eosio-ship-trace-reader/transport/redis_common"

	_redis "github.com/go-redis/redis/v8"
)

type Writer struct {
	pipeline _redis.Pipeliner
	ctx      context.Context
	ns       Namespace
	max_len  int64
}

type Option func(*Writer)

func WithMaxLen(value int64) Option {
	return func(p *Writer) {
		p.max_len = value
	}
}

func WithNamespace(value Namespace) Option {
	return func(p *Writer) {
		p.ns = value
	}
}

func NewWriter(client *_redis.Client, options ...Option) *Writer {
	pub := &Writer{
		pipeline: client.Pipeline(),
		ctx:      client.Context(),
		max_len:  2000,
	}

	for _, opt := range options {
		opt(pub)
	}

	return pub
}

func (r *Writer) Write(channel transport.Channel, payload []byte) error {
	args := &_redis.XAddArgs{
		Stream: r.ns.NewKey(channel).String(),
		ID:     "*",
		MaxLen: r.max_len,
		Values: payload,
	}

	return r.pipeline.XAdd(r.ctx, args).Err()
}

func (r *Writer) Flush() error {
	_, err := r.pipeline.Exec(r.ctx)
	return err
}
