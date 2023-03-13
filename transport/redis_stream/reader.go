package redis_stream

import (
	"context"
	"sync"
	"time"

	"eosio-ship-trace-reader/transport"
	. "eosio-ship-trace-reader/transport/redis_common"

	"github.com/go-redis/redis/v8"
)

type Reader struct {
	client   *redis.Client
	ctx      context.Context
	mu       sync.RWMutex
	timeout  time.Duration
	channels map[string]chan []byte
	ns       Namespace
}

type ReaderOption func(*Reader)

func WithTimeout(value time.Duration) ReaderOption {
	return func(s *Reader) {
		s.timeout = value
	}
}

func NewReader(client *redis.Client, ns Namespace, options ...ReaderOption) *Reader {
	sub := &Reader{
		client:   client,
		ctx:      client.Context(),
		channels: make(map[string]chan []byte),
		timeout:  time.Millisecond * 200,
		ns:       ns,
	}

	for _, opt := range options {
		opt(sub)
	}

	go sub.worker()

	return sub
}

// forward forwards a message to the channel.
// as writes to a unbuffered channel will block until it's read.
// We run select on it and discard the message if no read happends during timeout
func forward(msg redis.Message, ch chan<- []byte, timeout time.Duration) {
	select {
	case <-time.After(timeout):
	case ch <- []byte(msg.Payload):
	}
}

// worker reads messages from redis pubsub and forwards them to
// correct channels.
func (s *Reader) worker() {
	for name, chan := range s.channels {
		// Route message to correct channel.
		s.mu.RLock()
		if ch, ok := s.channels[msg.Channel]; ok {
			go forward(*msg, ch, s.timeout)
		}
		s.mu.RUnlock()
	}
}

func (r *Reader) Read(channel transport.Channel) ([]byte, error) {
	var err error

	r.client.XRange(r.ctx, channel.String(), "~", "~")

	key := s.ns.NewKey(channel).String()
	s.mu.RLock()
	ch, ok := s.channels[key]
	s.mu.RUnlock()

	return <-ch, nil
}
