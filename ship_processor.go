
package main

import (
    "log"
    "encoding/json"
    "github.com/eoscanada/eos-go/ship"
)

func processBlock(block *ship.GetBlocksResultV0) {

    if block.ThisBlock.BlockNum % 100 == 0 {
        log.Printf("Current: %d, Head: %d\n", block.ThisBlock.BlockNum, block.Head.BlockNum)
    }
}

func processTraces(traces []*ship.TransactionTraceV0) {

    for _, trace := range traces {

        payload, err := json.Marshal(trace)
        if err == nil {
            channel := RedisKey("transactions")
            if err := RedisPublish(channel, payload).Err(); err != nil {
                log.Printf("Failed to post to channel '%s': %s", channel, err)
            }
        } else {
            log.Println("Failed to encode transaction:", err)
        }

        // Actions
        for _, actionTraceVar := range trace.ActionTraces {
            trace := actionTraceVar.Impl.(*ship.ActionTraceV0)

            act := ActionTrace{
                Receiver: trace.Receiver,
                Contract: trace.Act.Account,
            	Action: trace.Act.Name,
            }

            abi, err := GetAbi(trace.Act.Account)
            if err == nil {
                v, err := DecodeAction(abi, trace.Act.Data, trace.Act.Name)
                if err != nil {
                    log.Print(err)
                }
                act.Data = v
            } else {
                log.Printf("Failed to get abi for contract %s: %s\n", trace.Act.Account, err)
            }


            payload, err := json.Marshal(act)
            if err != nil {
                log.Println("Failed to encode action:", err)
                continue
            }

            channels := []string{
                RedisKey("actions"),
                RedisKey(string(act.Contract), "actions"),
                RedisKey(string(act.Contract), "actions", string(act.Action)),
            }

            for _, channel := range channels {
                if err := RedisPublish(channel, payload).Err(); err != nil {
                    log.Printf("Failed to post to channel '%s': %s", channel, err)
                }
            }
        }
    }
}
