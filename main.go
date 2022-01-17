
package main

import (
    "os"
    "os/signal"
    "context"
    "log"
    "time"
    eos "github.com/eoscanada/eos-go"
    shipclient "github.com/eosswedenorg-go/eos-ship-client"
)

// ---------------------------
//  Global variables
// ---------------------------

var config Config

var chainInfo *eos.InfoResp

var shClient *shipclient.ShipClient

var eosClient *eos.API
var eosClientCtx = context.Background()

func run() {

    // Create done and interrupt channels.
    done := make(chan bool)
    interrupt := make(chan os.Signal, 1)

    // Register interrupt channel to receive interrupt messages
    signal.Notify(interrupt, os.Interrupt)

    // Spawn message read loop in another thread.
    go func() {
        for {
            err := shClient.Read()
            if err != nil {
                log.Print(err.Error())

                // Bail out on socket read error.
                if err.Type == shipclient.ErrSockRead {
                    break
                }
            }
        }

        shClient.Close()

        // Reader exited. signal that we are done.
        done <- true
    }()

    // Enter event loop in main thread
    for {
        select {
        case <-interrupt:
            log.Println("Interrupt, closing")

            // Cleanly close the connection by sending a close message and then
            // waiting (with timeout) for the server to close the connection.
            shClient.SendCloseMessage()

            select {
                case <-done: log.Println("Closed")
                case <-time.After(time.Second * 10): log.Println("Timeout");
            }
            return
        case <-done:
            log.Println("Closed")
            return
        }
    }
}

func main() {

    var err error

    // Parse config
    config, err = LoadConfig("config.json")
    if err != nil {
        log.Println(err)
        return
    }

    // Connect to redis
    RedisConnect(config.Redis.Addr, config.Redis.Password, config.Redis.DB)

    // Connect client and get chain info.
    eosClient = eos.New(config.Api)
    chainInfo, err = eosClient.GetInfo(eosClientCtx)
    if err != nil {
        log.Println("Failed to get info:", err)
        return
    }

    redisPrefix += chainInfo.ChainID.String() + "."

    if config.StartBlockNum == NULL_BLOCK_NUMBER {

        if config.IrreversibleOnly {
            config.StartBlockNum = uint32(chainInfo.LastIrreversibleBlockNum)
        } else {
            config.StartBlockNum = uint32(chainInfo.HeadBlockNum)
        }
    }

    // Construct ship client
    shClient = shipclient.NewClient(config.StartBlockNum, config.EndBlockNum, config.IrreversibleOnly)
    shClient.BlockHandler = processBlock
    shClient.TraceHandler = processTraces

    err = shClient.Connect(config.ShipApi)
    if err != nil {
        log.Println(err)
        return
    }

    err = shClient.SendBlocksRequest()
    if err != nil {
        log.Println(err)
        return
    }

    log.Printf("Start: %d, End: %d", shClient.StartBlock, shClient.EndBlock)

    // Run the application
    run()
}
