package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"eosio-ship-trace-reader/redis"
	"github.com/eoscanada/eos-go/ship"
)

func processBlock(block *ship.GetBlocksResultV0) {
	if block.ThisBlock.BlockNum%100 == 0 {
		log.Infof("Current: %d, Head: %d\n", block.ThisBlock.BlockNum, block.Head.BlockNum)
	}
}

func processTraces(traces []*ship.TransactionTraceV0) {
	for _, trace := range traces {

		payload, err := json.Marshal(trace)
		if err == nil {
			channel := redis.Key("transactions")
			if err := redis.Publish(channel, payload).Err(); err != nil {
				log.WithError(err).Errorf("Failed to post to channel '%s'", channel)
			}
		} else {
			log.WithError(err).Warn("Failed to encode transaction")
		}

		// Actions
		for _, actionTraceVar := range trace.ActionTraces {
			trace := actionTraceVar.Impl.(*ship.ActionTraceV0)

			act := ActionTrace{
				Receiver: trace.Receiver,
				Contract: trace.Act.Account,
				Action:   trace.Act.Name,
			}

			abi, err := GetAbi(trace.Act.Account)
			if err == nil {
				v, err := DecodeAction(abi, trace.Act.Data, trace.Act.Name)
				if err != nil {
					log.WithError(err).Warn("Failed to decode action")
				}
				act.Data = v
			} else {
				log.WithError(err).Errorf("Failed to get abi for contract %s", trace.Act.Account)
			}

			payload, err := json.Marshal(act)
			if err != nil {
				log.WithError(err).Error("Failed to encode action")
				continue
			}

			channels := []string{
				redis.Key("actions"),
				redis.Key(string(act.Contract), "actions"),
				redis.Key(string(act.Contract), "actions", string(act.Action)),
			}

			for _, channel := range channels {
				if err := redis.RegisterPublish(channel, payload).Err(); err != nil {
					log.WithError(err).Errorf("Failed to post to channel '%s'", channel)
				}
			}
		}
	}

	_, err := redis.Send()
	if err != nil {
		log.WithError(err).Error("Failed to send redis")
	}
}
