package api

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/eosswedenorg/thalos/api/message"
)

type handler func([]byte)

// Client reads and decodes messages from a reader and posts them to a go channel
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

// Helper method to post a message to a channel with timeout.
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
			// Don't report EOF as an error because it is used
			// by readers to signal an graceful end of input.
			if err != io.EOF {
				c.post(err)
			}
			return
		}

		h(payload)
	}
}

// Helper method to decode a message and post and error on the channel if it fails.
// Returns true if successful. False otherwise
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

func (c *Client) sub(channel Channel) error {
	var h handler

	switch channel.Type() {
	case RollbackChannel.Type():
		h = c.rollbackHandler
	case TransactionChannel.Type():
		h = c.transactionHandler
	case HeartbeatChannel.Type():
		h = c.hbHandler
	case ActionChannel{}.Channel().Type():
		h = c.actHandler
	case TableDeltaChannel{}.Channel().Type():
		h = c.tableDeltaHandler
	default:
		return fmt.Errorf("invalid channel type. %s", channel.Type())
	}

	// Start a worker for this channel.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.worker(channel, h)
	}()

	return nil
}

func (c *Client) Subscribe(channels ...Channel) error {
	for _, ch := range channels {
		if err := c.sub(ch); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Run() {
	// Just wait for workers to complete.
	c.wg.Wait()
}

func (c *Client) Close() error {
	err := c.reader.Close()
	// Wait for all goroutines to finish before closing channel.
	c.wg.Wait()
	close(c.channel)
	return err
}
