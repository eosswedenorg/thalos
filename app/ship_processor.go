package app

import (
	"encoding/hex"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	"github.com/eosswedenorg/thalos/app/abi"

	log "github.com/sirupsen/logrus"

	"github.com/eoscanada/eos-go/ship"
	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

// logDecoratedEncoder decorates a message.Encoder and logs any error.
func logDecoratedEncoder(encoder message.Encoder) message.Encoder {
	return func(v interface{}) ([]byte, error) {
		payload, err := encoder(v)
		if err != nil {
			log.WithError(err).
				WithField("v", v).
				Warn("Failed to encode message")
		}
		return payload, err
	}
}

type ShipProcessor struct {
	abi        *abi.AbiManager
	writer     api.Writer
	shipStream *shipclient.Stream
	encode     message.Encoder
}

func SpawnProccessor(shipStream *shipclient.Stream, writer api.Writer, abi *abi.AbiManager, codec message.Codec) *ShipProcessor {
	processor := &ShipProcessor{
		abi:        abi,
		writer:     writer,
		shipStream: shipStream,
		encode:     logDecoratedEncoder(codec.Encoder),
	}

	// Attach handlers
	shipStream.BlockHandler = processor.processBlock

	// Needed because if nil, traces will not be included in the response from ship.
	shipStream.TraceHandler = func([]*ship.TransactionTraceV0) {}

	return processor
}

func (processor *ShipProcessor) queueMessage(channel api.Channel, payload []byte) bool {
	err := processor.writer.Write(channel, payload)
	if err != nil {
		log.WithError(err).Errorf("Failed to post to channel '%s'", channel)
		return false
	}
	return true
}

func (processor *ShipProcessor) encodeQueue(channel api.Channel, v interface{}) bool {
	if payload, err := processor.encode(v); err == nil {
		return processor.queueMessage(channel, payload)
	}
	return false
}

func (processor *ShipProcessor) processBlock(block *ship.GetBlocksResultV0) {
	if block.ThisBlock.BlockNum%100 == 0 {
		log.Infof("Current: %d, Head: %d", block.ThisBlock.BlockNum, block.Head.BlockNum)
	}

	if block.ThisBlock.BlockNum%10 == 0 {
		hb := message.HeartBeat{
			BlockNum:                 block.ThisBlock.BlockNum,
			LastIrreversibleBlockNum: block.LastIrreversible.BlockNum,
			HeadBlockNum:             block.Head.BlockNum,
		}

		processor.encodeQueue(api.HeartbeatChannel, hb)
	}

	// Process traces
	if block.Traces != nil && len(block.Traces.Elem) > 0 {
		for _, trace := range block.Traces.AsTransactionTracesV0() {

			processor.encodeQueue(api.TransactionChannel, trace)

			// Actions
			for _, actionTraceVar := range trace.ActionTraces {
				var act_trace *ship.ActionTraceV1

				if trace_v0, ok := actionTraceVar.Impl.(*ship.ActionTraceV0); ok {
					// convert to v1
					act_trace = &ship.ActionTraceV1{
						ActionOrdinal:        trace_v0.ActionOrdinal,
						CreatorActionOrdinal: trace_v0.CreatorActionOrdinal,
						Receipt:              trace_v0.Receipt,
						Receiver:             trace_v0.Receiver,
						Act:                  trace_v0.Act,
						ContextFree:          trace_v0.ContextFree,
						Elapsed:              trace_v0.Elapsed,
						Console:              trace_v0.Console,
						AccountRamDeltas:     trace_v0.AccountRamDeltas,
						Except:               trace_v0.Except,
						ErrorCode:            trace_v0.ErrorCode,
						ReturnValue:          []byte{},
					}
				} else {
					act_trace = actionTraceVar.Impl.(*ship.ActionTraceV1)
				}

				act := message.ActionTrace{
					TxID:      trace.ID.String(),
					BlockNum:  block.Block.BlockNumber(),
					Timestamp: block.Block.Timestamp.Time.UTC(),
					Name:      act_trace.Act.Name.String(),
					Contract:  act_trace.Act.Account.String(),
					Receiver:  act_trace.Receiver.String(),
					HexData:   hex.EncodeToString(act_trace.Act.Data),
				}

				for _, auth := range act_trace.Act.Authorization {
					act.Authorization = append(act.Authorization, message.PermissionLevel{
						Actor:      auth.Actor.String(),
						Permission: auth.Permission.String(),
					})
				}

				ABI, err := processor.abi.GetAbi(act_trace.Act.Account)
				if err == nil {
					data, err := ABI.DecodeAction(act_trace.Act.Data, act_trace.Act.Name)
					if err != nil {
						log.WithError(err).Warn("Failed to decode action")
					}
					act.Data = data
				} else {
					log.WithError(err).Errorf("Failed to get abi for contract %s", act_trace.Act.Account)
				}

				payload, err := processor.encode(act)
				if err != nil {
					continue
				}

				channels := []api.Channel{
					api.ActionChannel{}.Channel(),
					api.ActionChannel{Name: act.Name}.Channel(),
					api.ActionChannel{Contract: act.Contract}.Channel(),
					api.ActionChannel{Name: act.Name, Contract: act.Contract}.Channel(),
				}

				for _, channel := range channels {
					processor.queueMessage(channel, payload)
				}
			}
		}
	}

	err := processor.writer.Flush()
	if err != nil {
		log.WithError(err).Error("Failed to send messages")
	}
}

func (processor *ShipProcessor) processTraces([]*ship.TransactionTraceV0) {
}

func (processor *ShipProcessor) Close() error {
	return processor.writer.Close()
}
