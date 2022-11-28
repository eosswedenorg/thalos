
package main

import (
    "log"
    "encoding/json"
    "github.com/eoscanada/eos-go/ship"
)

var block_num uint32

func processBlock(block *ship.GetBlocksResultV0) {

    block_num = block.ThisBlock.BlockNum

    if block_num % 100 == 0 {
        log.Printf("Current: %d, Head: %d\n", block_num, block.Head.BlockNum)
    }
}

func processTraces(traces []*ship.TransactionTraceV0) {

    for _, trace := range traces {

        payload, err := json.Marshal(trace)
        if err == nil {
            if err := transporter.Send("transactions", block_num, payload); err != nil {
                log.Println(err)
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
                "actions",
                string(act.Contract) + ".actions",
                string(act.Contract) + ".actions." + string(act.Action),
            }

            for _, channel := range channels {
                if err := transporter.Send(channel, block_num, payload); err != nil {
                    log.Println(err)
                }
            }
        }
    }

    err := transporter.Commit()
    if err != nil {
        log.Println("Failed to flush queue", err)
    }
}
