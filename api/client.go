package api

import (
	"sync"

	"github.com/eosswedenorg/thalos/api/message"
)

type handler func([]byte)

// Client reads and decodes messages from a reader and provides callback functions.
type Client struct {
	reader  Reader
	decoder message.Decoder

	wg sync.WaitGroup

	OnError     func(error)
	OnAction    func(message.ActionTrace)
	OnHeartbeat func(message.HeartBeat)
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

func (c *Client) Subscribe(channel Channel) {
	var handler handler

	if HeartbeatChannel.Is(channel) {
		handler = c.hbHandler
	} else {
		handler = c.actHandler
	}

	// Start a worker for this channel.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.worker(channel, handler)
	}()
}

func (c *Client) Run() {
	c.wg.Wait()
}

func (c *Client) Close() error {
	return c.reader.Close()
}
