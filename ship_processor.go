package main

import (
	"encoding/hex"
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"eosio-ship-trace-reader/internal/redis"
	"github.com/eoscanada/eos-go/ship"
)

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

func queueMessage(channel string, payload []byte) bool {
	err := redis.RegisterPublish(channel, payload).Err()
	if err != nil {
		log.WithError(err).Errorf("Failed to post to channel '%s'", channel)
		return false
	}
	return true
}

func encodeQueue(key string, v interface{}) bool {
	if payload, ok := encodeMessage(v); ok {
		channel := redis.Key(key)
		if queueMessage(channel, payload) {
			return true
		}
	}
	return false
}

func processBlock(block *ship.GetBlocksResultV0) {
	if block.ThisBlock.BlockNum%100 == 0 {
		log.Infof("Current: %d, Head: %d\n", block.ThisBlock.BlockNum, block.Head.BlockNum)
	}

	if block.ThisBlock.BlockNum%10 == 0 {
		hb := HearthBeat{
			BlockNum:                 block.ThisBlock.BlockNum,
			LastIrreversibleBlockNum: block.LastIrreversible.BlockNum,
			HeadBlockNum:             block.Head.BlockNum,
		}

		encodeQueue("hearthbeat", hb)

		_, err := redis.Send()
		if err != nil {
			log.WithError(err).Error("Failed to send redis")
		}
	}
}

func processTraces(traces []*ship.TransactionTraceV0) {
	for _, trace := range traces {

		encodeQueue("transactions", trace)

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

			abi, err := GetAbi(act_trace.Act.Account)
			if err == nil {
				v, err := DecodeAction(abi, act_trace.Act.Data, act_trace.Act.Name)
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

			channels := []string{
				redis.Key("actions"),
				redis.Key(string(act.Contract), "actions"),
				redis.Key(string(act.Contract), "actions", string(act.Action)),
			}

			for _, channel := range channels {
				queueMessage(channel, payload)
			}
		}
	}

	_, err := redis.Send()
	if err != nil {
		log.WithError(err).Error("Failed to send redis")
	}
}
