package app

import (
	"encoding/hex"

	log "github.com/sirupsen/logrus"

	"eosio-ship-trace-reader/abi"
	"eosio-ship-trace-reader/transport"
	"eosio-ship-trace-reader/transport/message"
	"github.com/eoscanada/eos-go/ship"
	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

type ShipProcessor struct {
	abi       *abi.AbiManager
	publisher transport.Publisher
	shClient  *shipclient.ShipClient
}

func SpawnProccessor(shClient *shipclient.ShipClient, publisher transport.Publisher, abi *abi.AbiManager) {
	processor := &ShipProcessor{
		abi:       abi,
		publisher: publisher,
		shClient:  shClient,
	}

	// Attach handlers
	shClient.BlockHandler = processor.processBlock
	shClient.TraceHandler = processor.processTraces
}

func (processor *ShipProcessor) queueMessage(channel transport.ChannelInterface, payload []byte) bool {
	err := processor.publisher.Publish(channel, payload)
	if err != nil {
		log.WithError(err).Errorf("Failed to post to channel '%s'", channel)
		return false
	}
	return true
}

func (processor *ShipProcessor) encodeQueue(channel transport.ChannelInterface, v interface{}) bool {
	if payload, ok := message.Encode(v); ok {
		if processor.queueMessage(channel, payload) {
			return true
		}
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

		err := processor.publisher.Flush()
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
				Receiver: act_trace.Receiver.String(),
				Contract: act_trace.Act.Account.String(),
				Action:   act_trace.Act.Name.String(),
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

			payload, ok := message.Encode(act)
			if !ok {
				continue
			}

			channels := []transport.ChannelInterface{
				transport.ActionChannel{},
				transport.ActionChannel{Action: string(act.Action)},
				transport.ActionChannel{Contract: string(act.Contract)},
				transport.ActionChannel{Action: string(act.Action), Contract: string(act.Contract)},
			}

			for _, channel := range channels {
				processor.queueMessage(channel, payload)
			}
		}
	}

	err := processor.publisher.Flush()
	if err != nil {
		log.WithError(err).Error("Failed to send messages")
	}
}
