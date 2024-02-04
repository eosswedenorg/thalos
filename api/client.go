package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/eosswedenorg/thalos/api/message"
)

type handler func([]byte)

// Client reads and decodes messages from a reader and posts thems to a go channel
type Client struct {
	reader  Reader
	decoder message.Decoder

	// waitgroup for worker threads.
	wg sync.WaitGroup

	// Channel for messages and errors
	channel chan any
}

func NewClient(reader Reader, decoder message.Decoder) *Client {
	return &Client{
		reader:  reader,
		decoder: decoder,
		channel: make(chan any),
	}
}

func (c *Client) Channel() <-chan any {
	return c.channel
}

func (c *Client) post(msg any) {
	select {
	case <-time.After(time.Second):
	case c.channel <- msg:
	}
}

func (c *Client) worker(channel Channel, h handler) {
	for {
		payload, err := c.reader.Read(channel)
		if err != nil {
			c.post(err)
			return
		}

		h(payload)
	}
}

// Helper method to decode a message and post and error on the channel if it fails.
// Returns true if successfull. false otherwise
func (c *Client) decode(payload []byte, msg any) bool {
	if err := c.decoder(payload, msg); err != nil {
		c.post(err)
		return false
	}
	return true
}

// Rollback handler
func (c *Client) rollbackHandler(payload []byte) {
	var rb message.RollbackMessage
	if ok := c.decode(payload, &rb); ok {
		c.post(rb)
	}
}

// Transaction handler
func (c *Client) transactionHandler(payload []byte) {
	var trans message.TransactionTrace
	if ok := c.decode(payload, &trans); ok {
		c.post(trans)
	}
}

// Action handler
func (c *Client) actHandler(payload []byte) {
	var act message.ActionTrace
	if ok := c.decode(payload, &act); ok {
		c.post(act)
	}
}

// TableDelta handler
func (c *Client) tableDeltaHandler(payload []byte) {
	td := message.TableDelta{}
	if ok := c.decode(payload, &td); ok {
		c.post(td)
	}
}

// HeartBeat handler
func (c *Client) hbHandler(payload []byte) {
	var hb message.HeartBeat
	if ok := c.decode(payload, &hb); ok {
		c.post(hb)
	}
}

func (c *Client) Subscribe(channel Channel) error {
	var handler handler

	switch channel.Type() {
	case RollbackChannel.Type():
		handler = c.rollbackHandler
	case TransactionChannel.Type():
		handler = c.transactionHandler
	case HeartbeatChannel.Type():
		handler = c.hbHandler
	case ActionChannel{}.Channel().Type():
		handler = c.actHandler
	case TableDeltaChannel{}.Channel().Type():
		handler = c.tableDeltaHandler
	default:
		return fmt.Errorf("invalid channel type. %s", channel.Type())
	}

	// Start a worker for this channel.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.worker(channel, handler)
	}()

	return nil
}

func (c *Client) Run() {
	// Just wait for workers to complete.
	c.wg.Wait()
}

func (c *Client) Close() error {
	err := c.reader.Close()
	// Wait for all goroutines before closing channel.
	c.wg.Wait()
	close(c.channel)
	return err
}
