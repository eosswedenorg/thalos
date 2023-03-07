package app

import (
	"encoding/hex"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"eosio-ship-trace-reader/abi"
	"eosio-ship-trace-reader/transport"
	"eosio-ship-trace-reader/transport/message"
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
	abi      *abi.AbiManager
	writer   transport.Writer
	shClient *shipclient.Client
	encode   message.Encoder
}

func SpawnProccessor(shClient *shipclient.Client, writer transport.Writer, abi *abi.AbiManager) {
	processor := &ShipProcessor{
		abi:      abi,
		writer:   writer,
		shClient: shClient,
		encode:   logDecoratedEncoder(json.Marshal),
	}

	// Attach handlers
	shClient.BlockHandler = processor.processBlock
	shClient.TraceHandler = processor.processTraces
}

func (processor *ShipProcessor) queueMessage(channel transport.Channel, payload []byte) bool {
	err := processor.writer.Write(channel, payload)
	if err != nil {
		log.WithError(err).Errorf("Failed to post to channel '%s'", channel)
		return false
	}
	return true
}

func (processor *ShipProcessor) encodeQueue(channel transport.Channel, v interface{}) bool {
	if payload, err := processor.encode(v); err == nil {
		return processor.queueMessage(channel, payload)
	}
	return false
}

func (processor *ShipProcessor) processBlock(block *ship.GetBlocksResultV0) {
	if block.ThisBlock.BlockNum%100 == 0 {
		log.Infof("Current: %d, Head: %d\n", block.ThisBlock.BlockNum, block.Head.BlockNum)
	}

	if block.ThisBlock.BlockNum%10 == 0 {
		hb := message.HearthBeat{
			BlockNum:                 block.ThisBlock.BlockNum,
			LastIrreversibleBlockNum: block.LastIrreversible.BlockNum,
			HeadBlockNum:             block.Head.BlockNum,
		}

		processor.encodeQueue(transport.HeartbeatChannel, hb)

		err := processor.writer.Flush()
		if err != nil {
			log.WithError(err).Error("Failed to send messages")
		}
	}
}

func (processor *ShipProcessor) processTraces(traces []*ship.TransactionTraceV0) {
	for _, trace := range traces {

		processor.encodeQueue(transport.TransactionChannel, trace)

		// Actions
		for _, actionTraceVar := range trace.ActionTraces {
			act_trace := actionTraceVar.Impl.(*ship.ActionTraceV0)

			act := message.ActionTrace{
				TxID:     trace.ID.String(),
				Name:     act_trace.Act.Name.String(),
				Contract: act_trace.Act.Account.String(),
				Receiver: act_trace.Receiver.String(),
				HexData:  hex.EncodeToString(act_trace.Act.Data),
			}

			ABI, err := processor.abi.GetAbi(act_trace.Act.Account)
			if err == nil {
				v, err := abi.DecodeAction(ABI, act_trace.Act.Data, act_trace.Act.Name)
				if err != nil {
					log.WithError(err).Warn("Failed to decode action")
				}
				act.Data = v
			} else {
				log.WithError(err).Errorf("Failed to get abi for contract %s", act_trace.Act.Account)
			}

			payload, err := processor.encode(act)
			if err != nil {
				continue
			}

			channels := []transport.Channel{
				transport.Action{}.Channel(),
				transport.Action{Name: act.Name}.Channel(),
				transport.Action{Contract: act.Contract}.Channel(),
				transport.Action{Name: act.Name, Contract: act.Contract}.Channel(),
			}

			for _, channel := range channels {
				processor.queueMessage(channel, payload)
			}
		}
	}

	err := processor.writer.Flush()
	if err != nil {
		log.WithError(err).Error("Failed to send messages")
	}
}
