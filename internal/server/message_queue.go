package server

import (
	"errors"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	"github.com/eosswedenorg/thalos/internal/driver"
)

// MessageQueue takes care of message routing and encoding
type MessageQueue struct {
	// Writer to write messages to
	writer driver.Writer

	// Encoder to encode messages with
	encode message.Encoder
}

func NewMessageQueue(writer driver.Writer, encoder message.Encoder) MessageQueue {
	return MessageQueue{
		writer: writer,
		encode: encoder,
	}
}

func (mq MessageQueue) PostHeartbeat(hb message.HeartBeat) error {
	return mq.post(hb, api.HeartbeatChannel)
}

func (mq MessageQueue) PostRollback(rb message.RollbackMessage) error {
	return mq.post(rb, api.RollbackChannel)
}

func (mq MessageQueue) PostTransactionTrace(trace message.TransactionTrace) error {
	return mq.post(trace, api.TransactionChannel)
}

// Post a ActionTrace message to the queue
func (mq MessageQueue) PostAction(act message.ActionTrace) error {
	return mq.post(act,
		api.ActionChannel{}.Channel(),
		api.ActionChannel{Name: act.Name}.Channel(),
		api.ActionChannel{Contract: act.Contract}.Channel(),
		api.ActionChannel{Name: act.Name, Contract: act.Contract}.Channel(),
	)
}

func (mq MessageQueue) PostTableDelta(delta message.TableDelta) error {
	return mq.post(delta,
		api.TableDeltaChannel{}.Channel(),
		api.TableDeltaChannel{Name: delta.Name}.Channel(),
	)
}

func (mq MessageQueue) Flush() error {
	return mq.writer.Flush()
}

func (mq MessageQueue) Close() error {
	return mq.writer.Close()
}

func (mq MessageQueue) post(v interface{}, channels ...api.Channel) error {
	payload, err := mq.encode(v)
	if err == nil {
		for _, channel := range channels {
			if w_err := mq.writer.Write(channel, payload); err != nil {
				err = errors.Join(w_err)
			}
		}
	}
	return err
}
