package server

import (
	"bytes"

	"github.com/eosswedenorg/thalos/api/message"
	"github.com/eosswedenorg/thalos/internal/abi"
	"github.com/eosswedenorg/thalos/internal/driver"
	ship_helper "github.com/eosswedenorg/thalos/internal/ship"
	"github.com/eosswedenorg/thalos/internal/types"

	log "github.com/sirupsen/logrus"

	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
	"github.com/shufflingpixels/antelope-go/chain"
	"github.com/shufflingpixels/antelope-go/ship"
)

// A ShipProcessor will consume messages from a ship stream, convert the messages into
// thalos specific ones, encode them and finally post them to an api.Writer
type ShipProcessor struct {
	// The ship stream to process.
	shipStream *shipclient.Stream

	// Abi manager used for cacheing
	abi *abi.AbiManager

	queue MessageQueue

	// Function for saving state.
	saver StateSaver

	// Internal state
	state State

	// System contract ("eosio" per default)
	syscontract chain.Name

	// ABI Returned from SHIP
	shipABI *chain.Abi

	// Action blacklist
	blacklist types.Blacklist
}

// SpawnProcessor creates a new ShipProccessor that consumes the shipclient.Stream passed to it.
func SpawnProccessor(shipStream *shipclient.Stream, loader StateLoader, saver StateSaver, writer driver.Writer, abi *abi.AbiManager, codec message.Codec) *ShipProcessor {
	processor := &ShipProcessor{
		saver:       saver,
		abi:         abi,
		shipStream:  shipStream,
		syscontract: chain.N("eosio"),
		queue:       NewMessageQueue(writer, codec.Encoder),
	}

	loader(&processor.state)

	// Attach handlers
	shipStream.BlockHandler = processor.processBlock
	shipStream.InitHandler = processor.initHandler

	// Needed because if nil, traces/table deltas will not be included in the response from ship.
	shipStream.TraceHandler = func(*ship.TransactionTraceArray) {}
	shipStream.TableDeltaHandler = func(*ship.TableDeltaArray) {}

	return processor
}

func (processor *ShipProcessor) SetBlacklist(list types.Blacklist) {
	processor.blacklist = list
}

func (processor *ShipProcessor) initHandler(abi *chain.Abi) {
	processor.shipABI = abi
}

// updateAbiFromAction updates the contract abi based on the ship.Action passed.
func (processor *ShipProcessor) updateAbiFromAction(act *chain.Action) error {
	set_abi := struct {
		Account chain.Name
		Abi     chain.Bytes
	}{}

	if err := act.DecodeInto(&set_abi); err != nil {
		return err
	}

	abi := chain.Abi{}
	decoder := chain.NewDecoder(bytes.NewReader(set_abi.Abi))
	if err := decoder.Decode(&abi); err != nil {
		return err
	}
	return processor.abi.SetAbi(set_abi.Account, &abi)
}

// Get the current block.
func (processor *ShipProcessor) GetCurrentBlock() uint32 {
	return processor.state.CurrentBlock
}

func (processor *ShipProcessor) processTransactionTrace(log *log.Entry, blockNumber uint32, block *ship.SignedBlock, trace *ship.TransactionTraceV0) {
	logger := log.WithField("type", "trace").WithField("tx_id", trace.ID.String()).Dup()

	timestamp := block.BlockHeader.Timestamp.Time().UTC()

	transaction := message.TransactionTrace{
		ID:            trace.ID.String(),
		BlockNum:      blockNumber,
		Timestamp:     timestamp,
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

		actionTrace := ship_helper.ToActionTraceV1(actionTraceVar)
		actMsg := processor.proccessActionTrace(logger, actionTrace)
		if actMsg != nil {
			actMsg.TxID = trace.ID.String()
			actMsg.BlockNum = blockNumber
			actMsg.Timestamp = timestamp

			processor.queue.PostAction(*actMsg)

			transaction.ActionTraces = append(transaction.ActionTraces, *actMsg)
		}
	}

	if err := processor.queue.PostTransactionTrace(transaction); err != nil {
		logger.WithError(err).Error("Failed to post transaction trace")
	}
}

func (processor *ShipProcessor) proccessActionTrace(logger *log.Entry, trace *ship.ActionTraceV1) *message.ActionTrace {
	// Check if actions updates an abi.
	if trace.Act.Account == processor.syscontract && trace.Act.Name == chain.N("setabi") {

		logger.WithFields(log.Fields{
			"contract": trace.Act.Account,
			"action":   trace.Act.Name,
		}).Debug("Update contract ABI")

		err := processor.updateAbiFromAction(&trace.Act)
		if err != nil {
			logger.WithError(err).Warn("Failed to update abi")
		}
	}

	// Check blacklist if we should skip this action
	if !processor.blacklist.IsAllowed(trace.Act.Account.String(), trace.Act.Name.String()) {
		logger.WithFields(log.Fields{
			"contract": trace.Act.Account,
			"action":   trace.Act.Name,
		}).Debug("Found in blacklist, skipping")
		return nil
	}

	act := &message.ActionTrace{
		Name:          trace.Act.Name.String(),
		Contract:      trace.Act.Account.String(),
		Receiver:      trace.Receiver.String(),
		FirstReceiver: trace.Act.Account.String() == trace.Receiver.String(),
	}

	if trace.Receipt != nil {
		receipt := trace.Receipt.V0
		act.Receipt = &message.ActionReceipt{
			Receiver:       receipt.Receiver.String(),
			ActDigest:      receipt.ActDigest.String(),
			GlobalSequence: receipt.GlobalSequence,
			RecvSequence:   receipt.RecvSequence,
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

	for _, auth := range trace.Act.Authorization {
		act.Authorization = append(act.Authorization, message.PermissionLevel{
			Actor:      auth.Actor.String(),
			Permission: auth.Permission.String(),
		})
	}

	logger.WithFields(log.Fields{
		"contract": trace.Act.Account,
		"action":   trace.Act.Name,
	}).Debug("Reading contract ABI")

	ABI, err := processor.abi.GetAbi(trace.Act.Account)
	if err == nil {
		if act.Data, err = trace.Act.Decode(ABI); err != nil {
			logger.WithFields(log.Fields{
				"contract": trace.Act.Account,
				"action":   trace.Act.Name,
			}).WithError(err).Warn("Failed to decode action")
		}
	} else {
		logger.WithField("contract", trace.Act.Account).
			WithError(err).Error("Failed to get abi for contract")
	}

	return act
}

func (processor *ShipProcessor) proccessDeltaRows(logger *log.Entry, table_name string, rows []ship.Row) []message.TableDeltaRow {
	out := []message.TableDeltaRow{}
	for _, row := range rows {

		msg := message.TableDeltaRow{
			Present: row.Present,
			RawData: row.Data,
		}

		if processor.shipABI != nil {

			v, err := processor.shipABI.Decode(bytes.NewReader(row.Data), table_name)
			if err == nil {
				v, err := ship_helper.ParseTableDeltaData(v)
				if err == nil {
					msg.Data = v
				} else {
					logger.WithError(err).Error("Failed to parse table delta data")
				}
			} else {
				logger.Error("Failed to decode table delta")
			}
		} else {
			logger.Warn("No SHIP ABI present")
		}
		out = append(out, msg)
	}
	return out
}

// Callback function called by shipclient.Stream when a new block arrives.
func (processor *ShipProcessor) processBlock(blockResult *ship.GetBlocksResultV0) {
	block := ship.SignedBlock{}
	blockResult.Block.Unpack(&block)
	timestamp := block.BlockHeader.Timestamp.Time().UTC()
	blockNumber := blockResult.ThisBlock.BlockNum

	// Check to see if we have a microfork and post a message to
	// the rollback channel in that case.
	if processor.state.CurrentBlock > 0 && blockNumber < processor.state.CurrentBlock {

		msg := message.RollbackMessage{
			OldBlockNum: processor.state.CurrentBlock,
			NewBlockNum: blockResult.ThisBlock.BlockNum,
		}
		log.WithField("old_block", msg.OldBlockNum).
			WithField("new_block", msg.NewBlockNum).
			Warn("Fork detected, old_block is greater than new_block")

		if err := processor.queue.PostRollback(msg); err != nil {
			log.WithError(err).Error("Failed to write rollback message")
		}
	}

	processor.state.CurrentBlock = blockNumber

	if blockResult.ThisBlock.BlockNum%100 == 0 {
		log.Infof("Current: %d, Head: %d", processor.state.CurrentBlock, blockResult.Head.BlockNum)
	}

	if blockResult.ThisBlock.BlockNum%10 == 0 {
		hb := message.HeartBeat{
			BlockNum:                 blockNumber,
			LastIrreversibleBlockNum: blockResult.LastIrreversible.BlockNum,
			HeadBlockNum:             blockResult.Head.BlockNum,
		}
		if err := processor.queue.PostHeartbeat(hb); err != nil {
			log.WithError(err).Error("Failed to write heartbeat message")
		}
	}

	mainLogger := log.WithField("block", blockNumber).Dup()

	// Process traces
	if blockResult.Traces != nil {
		unpacked := []ship.TransactionTrace{}
		if err := blockResult.Traces.Unpack(&unpacked); err != nil {
			mainLogger.WithError(err).Error("Failed to unpack transaction traces")
		} else {
			for _, trace := range unpacked {
				processor.processTransactionTrace(mainLogger, blockNumber, &block, trace.V0)
			}
		}
	}

	// Process deltas
	if blockResult.Deltas != nil {
		deltas := []ship.TableDelta{}
		if err := blockResult.Deltas.Unpack(&deltas); err != nil {
			mainLogger.WithError(err).Error("Failed to unpack table deltas")
		} else {
			logger := mainLogger.WithField("type", "table_delta").Dup()
			for _, delta := range deltas {

				msg := message.TableDelta{
					BlockNum:  blockNumber,
					Timestamp: timestamp,
					Name:      delta.V0.Name,
					Rows:      processor.proccessDeltaRows(logger, delta.V0.Name, delta.V0.Rows),
				}

				if err := processor.queue.PostTableDelta(msg); err != nil {
					logger.WithError(err).Error("Failed to post table delta message")
				}
			}
		}
	}

	err := processor.queue.Flush()
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
	return processor.queue.Close()
}
