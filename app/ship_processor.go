package app

import (
	"encoding/hex"
	"encoding/json"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	"github.com/eosswedenorg/thalos/app/abi"
	"github.com/eosswedenorg/thalos/app/driver"

	log "github.com/sirupsen/logrus"

	"github.com/eoscanada/eos-go"
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

// A ShipProcessor will consume messages from a ship stream, convert the messages into
// thalos specfic ones, encode them and finally post them to an api.Writer
type ShipProcessor struct {
	// The ship stream to process.
	shipStream *shipclient.Stream

	// Abi manager used for cacheing
	abi *abi.AbiManager

	// Writer to send messages to.
	writer driver.Writer

	// Encoder used to encode messages
	encode message.Encoder

	// Function for saving state.
	saver StateSaver

	// Internal state
	state State

	// System contract ("eosio" per default)
	syscontract eos.AccountName
}

// SpawnProcessor creates a new ShipProccessor that consumes the shipclient.Stream passed to it.
func SpawnProccessor(shipStream *shipclient.Stream, loader StateLoader, saver StateSaver, writer driver.Writer, abi *abi.AbiManager, codec message.Codec) *ShipProcessor {
	processor := &ShipProcessor{
		saver:       saver,
		abi:         abi,
		writer:      writer,
		shipStream:  shipStream,
		encode:      logDecoratedEncoder(codec.Encoder),
		syscontract: eos.AccountName("eosio"),
	}

	loader(&processor.state)

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

func decode(abi *eos.ABI, act *ship.Action, v any) error {
	jsondata, err := abi.DecodeAction(act.Data, act.Name)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsondata, v)
}

// updateAbiFromAction updates the contract abi based on the ship.Action passed.
func (processor *ShipProcessor) updateAbiFromAction(act *ship.Action) error {
	ABI, err := processor.abi.GetAbi(processor.syscontract)
	if err != nil {
		return err
	}

	set_abi := struct {
		Abi     string
		Account eos.AccountName
	}{}
	if err := decode(ABI, act, &set_abi); err != nil {
		return err
	}

	binary_abi, err := hex.DecodeString(set_abi.Abi)
	if err != nil {
		return err
	}

	contract_abi := eos.ABI{}
	if err = eos.UnmarshalBinary(binary_abi, &contract_abi); err != nil {
		return err
	}

	return processor.abi.SetAbi(set_abi.Account, &contract_abi)
}

// Get the current block.
func (processor *ShipProcessor) GetCurrentBlock() uint32 {
	return processor.state.CurrentBlock
}

// Callback function called by shipclient.Stream when a new block arrives.
func (processor *ShipProcessor) processBlock(block *ship.GetBlocksResultV0) {
	processor.state.CurrentBlock = block.ThisBlock.BlockNum

	if block.ThisBlock.BlockNum%100 == 0 {
		log.Infof("Current: %d, Head: %d", processor.state.CurrentBlock, block.Head.BlockNum)
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

			logger := log.WithField("tx_id", trace.ID.String()).Dup()

			transaction := message.TransactionTrace{
				ID:            trace.ID.String(),
				BlockNum:      block.Block.BlockNumber(),
				Timestamp:     block.Block.Timestamp.Time.UTC(),
				Status:        trace.Status.String(),
				CPUUsageUS:    trace.CPUUsageUS,
				NetUsage:      trace.NetUsage,
				NetUsageWords: uint32(trace.NetUsageWords),
				Elapsed:       int64(trace.Elapsed),
				Scheduled:     trace.Scheduled,
				Except:        trace.Except,
				Error:         trace.ErrorCode,
			}

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

				// Check if actions updates an abi.
				if act_trace.Act.Account == processor.syscontract && act_trace.Act.Name == eos.ActionName("setabi") {
					err := processor.updateAbiFromAction(act_trace.Act)
					if err != nil {
						logger.WithError(err).Warn("Failed to update abi")
					}
				}

				act := message.ActionTrace{
					TxID:      trace.ID.String(),
					BlockNum:  block.Block.BlockNumber(),
					Timestamp: block.Block.Timestamp.Time.UTC(),
					Name:      act_trace.Act.Name.String(),
					Contract:  act_trace.Act.Account.String(),
					Receiver:  act_trace.Receiver.String(),
				}

				if act_trace.Receipt != nil {
					receipt := act_trace.Receipt.Impl.(*ship.ActionReceiptV0)
					act.Receipt = &message.ActionReceipt{
						Receiver:       receipt.Receiver.String(),
						ActDigest:      receipt.ActDigest.String(),
						GlobalSequence: receipt.GlobalSequence,
						RecvSequence:   receipt.GlobalSequence,
						CodeSequence:   uint32(receipt.CodeSequence),
						ABISequence:    uint32(receipt.ABISequence),
					}

					for _, auth := range receipt.AuthSequence {
						act.Receipt.AuthSequence = append(act.Receipt.AuthSequence, message.AccountAuthSequence{
							Account:  auth.Account.String(),
							Sequence: auth.Sequence,
						})
					}
				}

				for _, auth := range act_trace.Act.Authorization {
					act.Authorization = append(act.Authorization, message.PermissionLevel{
						Actor:      auth.Actor.String(),
						Permission: auth.Permission.String(),
					})
				}

				ABI, err := processor.abi.GetAbi(act_trace.Act.Account)
				if err == nil {
					if err = decode(ABI, act_trace.Act, &act.Data); err != nil {
						logger.WithFields(log.Fields{
							"contract": act_trace.Act.Account,
							"action":   act_trace.Act.Name,
						}).WithError(err).Warn("Failed to decode action")
					}
				} else {
					logger.WithField("contract", act_trace.Act.Account).
						WithError(err).Error("Failed to get abi for contract")
				}

				payload, err := processor.encode(act)
				if err != nil {
					continue
				}

				transaction.ActionTraces = append(transaction.ActionTraces, act)

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

			processor.encodeQueue(api.TransactionChannel, transaction)
		}
	}

	err := processor.writer.Flush()
	if err != nil {
		log.WithError(err).Error("Failed to send messages")
	}

	err = processor.saver(processor.state)
	if err != nil {
		log.WithError(err).Error("Failed to save state")
	}
}

// Close closes the writer associated with the processor.
func (processor *ShipProcessor) Close() error {
	return processor.writer.Close()
}
