package api

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/eosswedenorg/thalos/api/message"
)

type handler func([]byte)

// Client reads and decodes messages from a reader and provides callback functions.
type Client struct {
	reader  Reader
	decoder message.Decoder

	// waitgroup for worker threads.
	wg sync.WaitGroup

	OnError      func(error)
	OnAction     func(message.ActionTrace)
	OnHeartbeat  func(message.HeartBeat)
	OnTableDelta func(message.TableDelta)
}

func NewClient(reader Reader, decoder message.Decoder) *Client {
	return &Client{
		reader:  reader,
		decoder: decoder,
	}
}

func (c *Client) worker(channel Channel, h handler) {
	for {
		payload, err := c.reader.Read(channel)
		if err != nil {
			if c.OnError != nil {
				c.OnError(err)
			}
			return
		}

		h(payload)
	}
}

// Action handler
func (c *Client) actHandler(payload []byte) {
	var act message.ActionTrace
	if err := c.decoder(payload, &act); err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		return
	}
	c.OnAction(act)
}

// TableDelta handler
func (c *Client) tableDeltaHandler(payload []byte) {
	td := message.TableDelta{}
	if err := c.decoder(payload, &td); err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		return
	}
	c.OnTableDelta(td)
}

// HeartBeat handler
func (c *Client) hbHandler(payload []byte) {
	var hb message.HeartBeat
	if err := c.decoder(payload, &hb); err != nil {
		if c.OnError != nil {
			c.OnError(err)
		}
		return
	}
	c.OnHeartbeat(hb)
}

func (c *Client) Subscribe(channel Channel) error {
	handlers := map[string]struct {
		handler  handler
		callback any
	}{
		HeartbeatChannel.Type():              {c.hbHandler, c.OnHeartbeat},
		ActionChannel{}.Channel().Type():     {c.actHandler, c.OnAction},
		TableDeltaChannel{}.Channel().Type(): {c.tableDeltaHandler, c.OnTableDelta},
	}

	h, ok := handlers[channel.Type()]

	if !ok {
		return fmt.Errorf("invalid channel type. %s", channel.Type())
	}

	if h.callback == nil || reflect.ValueOf(h.callback).IsNil() {
		return fmt.Errorf("please set an handler before calling Subscribe")
	}

	// Start a worker for this channel.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.worker(channel, h.handler)
	}()

	return nil
}

func (c *Client) Run() {
	// Just wait for workers to complete.
	c.wg.Wait()
}

func (c *Client) Close() error {
	return c.reader.Close()
}
