package transport

import (
	"encoding/json"

	"eosio-ship-trace-reader/transport/message"
)

type Client struct {
	reader  Reader
	decoder message.Decoder

	actChan chan message.ActionTrace
	hbChan  chan message.HearthBeat
	errChan chan error
}

func NewClient(reader Reader) *Client {
	return &Client{
		reader:  reader,
		decoder: json.Unmarshal,
		actChan: make(chan message.ActionTrace, 16),
		hbChan:  make(chan message.HearthBeat, 16),
		errChan: make(chan error, 16),
	}
}

func actWorker(decoder message.Decoder, out chan<- message.ActionTrace, reader Reader, channel Channel) {
	for {
		payload, err := reader.Read(channel)
		if err != nil {
			return
		}

		var act message.ActionTrace
		if err := decoder(payload, &act); err != nil {
			continue
		}
		out <- act
	}
}

func hbWorker(decoder message.Decoder, out chan<- message.HearthBeat, reader Reader, channel Channel) {
	for {
		payload, err := reader.Read(channel)
		if err != nil {
			return
		}

		var hb message.HearthBeat
		if err := decoder(payload, &hb); err != nil {
			continue
		}
		out <- hb
	}
}

func (c Client) Subscribe(channel Channel) {
	if HeartbeatChannel.Is(channel) {
		go hbWorker(c.decoder, c.hbChan, c.reader, channel)
	}

	go actWorker(c.decoder, c.actChan, c.reader, channel)
}

func (c Client) ActionTrace() <-chan message.ActionTrace {
	return c.actChan
}

func (c Client) Heartbeat() <-chan message.HearthBeat {
	return c.hbChan
}

func (c Client) Close() {
	close(c.actChan)
	close(c.hbChan)
}
