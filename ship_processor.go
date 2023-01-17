package main

import (
	"encoding/hex"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"eosio-ship-trace-reader/abi"
	"eosio-ship-trace-reader/transport"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ship"
)

type ShipReader struct {
	ns        transport.Namespace
	abi       *abi.AbiManager
	publisher transport.Publisher
}

func decodeAction(abi *eos.ABI, data []byte, actionName eos.ActionName) (interface{}, error) {
	var v interface{}

	bytes, err := abi.DecodeAction(data, actionName)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(bytes, &v)
	return v, err
}

func encodeMessage(v interface{}) ([]byte, bool) {
	payload, err := json.Marshal(v)
	if err != nil {
		log.WithError(err).
			WithField("v", v).
			Warn("Failed to encode message to json")
		return nil, false
	}
	return payload, true
}

func (reader *ShipReader) queueMessage(channel transport.ChannelInterface, payload []byte) bool {
	key := reader.ns.NewKey(channel)
	err := reader.publisher.Publish(key.String(), payload)
	if err != nil {
		log.WithError(err).Errorf("Failed to post to channel '%s'", key)
		return false
	}
	return true
}

func (reader *ShipReader) encodeQueue(channel transport.ChannelInterface, v interface{}) bool {
	if payload, ok := encodeMessage(v); ok {
		if reader.queueMessage(channel, payload) {
			return true
		}
	}
	return false
}

func (reader *ShipReader) processBlock(block *ship.GetBlocksResultV0) {
	if block.ThisBlock.BlockNum%100 == 0 {
		log.Infof("Current: %d, Head: %d\n", block.ThisBlock.BlockNum, block.Head.BlockNum)
	}

	if block.ThisBlock.BlockNum%10 == 0 {
		hb := HearthBeat{
			BlockNum:                 block.ThisBlock.BlockNum,
			LastIrreversibleBlockNum: block.LastIrreversible.BlockNum,
			HeadBlockNum:             block.Head.BlockNum,
		}

		reader.encodeQueue(transport.HeartbeatChannel, hb)

		err := reader.publisher.Flush()
		if err != nil {
			log.WithError(err).Error("Failed to send messages")
		}
	}
}

func (reader *ShipReader) processTraces(traces []*ship.TransactionTraceV0) {
	for _, trace := range traces {

		reader.encodeQueue(transport.TransactionChannel, trace)

		// Actions
		for _, actionTraceVar := range trace.ActionTraces {
			act_trace := actionTraceVar.Impl.(*ship.ActionTraceV0)

			act := ActionTrace{
				TxID:     trace.ID,
				Receiver: act_trace.Receiver,
				Contract: act_trace.Act.Account,
				Action:   act_trace.Act.Name,
				HexData:  hex.EncodeToString(act_trace.Act.Data),
			}

			abi, err := reader.abi.GetAbi(act_trace.Act.Account)
			if err == nil {
				v, err := decodeAction(abi, act_trace.Act.Data, act_trace.Act.Name)
				if err != nil {
					log.WithError(err).Warn("Failed to decode action")
				}
				act.Data = v
			} else {
				log.WithError(err).Errorf("Failed to get abi for contract %s", act_trace.Act.Account)
			}

			payload, ok := encodeMessage(act)
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
				reader.queueMessage(channel, payload)
			}
		}
	}

	err := reader.publisher.Flush()
	if err != nil {
		log.WithError(err).Error("Failed to send messages")
	}
}
