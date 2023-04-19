package redis_pubsub

import (
	"context"
	"sync"
	"time"

	"thalos/api"
	. "thalos/api/redis_common"

	"github.com/go-redis/redis/v8"
)

type Subscriber struct {
	client   *redis.Client
	sub      *redis.PubSub
	ctx      context.Context
	mu       sync.RWMutex
	timeout  time.Duration
	channels map[string]chan []byte
	ns       Namespace
}

type SubscriberOption func(*Subscriber)

func WithTimeout(value time.Duration) SubscriberOption {
	return func(s *Subscriber) {
		s.timeout = value
	}
}

func NewSubscriber(client *redis.Client, ns Namespace, options ...SubscriberOption) *Subscriber {
	sub := &Subscriber{
		client:   client,
		ctx:      client.Context(),
		sub:      client.PSubscribe(client.Context()),
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
func (s *Subscriber) worker() {
	for msg := range s.sub.Channel() {
		// Route message to correct channel.
		s.mu.RLock()
		if ch, ok := s.channels[msg.Channel]; ok {
			go forward(*msg, ch, s.timeout)
		}
		s.mu.RUnlock()
	}
}

func (s *Subscriber) Read(channel api.Channel) ([]byte, error) {
	var err error

	key := s.ns.NewKey(channel).String()
	s.mu.RLock()
	ch, ok := s.channels[key]
	s.mu.RUnlock()
	if !ok {
		// Channel does not exist in the map.
		// Subscribe and insert it.
		err = s.sub.Subscribe(s.ctx, key)
		if err != nil {
			return nil, err
		}

		// Guard race condition to map with mutex.
		s.mu.Lock()
		ch = make(chan []byte)
		s.channels[key] = ch
		s.mu.Unlock()
	}

	return <-ch, nil
}

func (s *Subscriber) Close() error {
	err := s.sub.Close()

	for _, ch := range s.channels {
		close(ch)
	}
	s.mu.Lock()
	s.channels = make(map[string]chan []byte)
	s.mu.Unlock()

	return err
}
