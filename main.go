
package main

import (
    "fmt"
    "os"
    "os/signal"
    "context"
    "log"
    "time"
    "github.com/pborman/getopt/v2"
    "github.com/eosswedenorg-go/pid"
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


// Reader states
const RS_CONNECT = 1
const RS_READ = 2

func readerLoop() {

    state := RS_CONNECT

    for {
        switch state {
        case RS_CONNECT :
            log.Printf("Connecting to ship at: %s", config.ShipApi)
            err := shClient.Connect(config.ShipApi)
            if err != nil {
                log.Println(err)
                log.Printf("Trying again in 5 seconds ....")
                time.Sleep(5 * time.Second)
                break;
            }

            err = shClient.SendBlocksRequest()
            if err != nil {
                log.Println(err)
                break
            }

            // Connected
            log.Printf("Connected, Start: %d, End: %d", shClient.StartBlock, shClient.EndBlock)
            state = RS_READ
        case RS_READ :
            err := shClient.Read()
            if err != nil {
                log.Print(err.Error())

                // Reconnect
                if err.Type == shipclient.ErrSockRead {
                    state = RS_CONNECT
                }
            }
        }
    }

    shClient.Close()
}

func run() {

    // Create done and interrupt channels.
    done := make(chan bool)
    interrupt := make(chan os.Signal, 1)

    // Register interrupt channel to receive interrupt messages
    signal.Notify(interrupt, os.Interrupt)

    // Spawn message read loop in another thread.
    go func() {
        readerLoop()

        // Reader exited. signal that we are done.
        done <- true
    }()

    // Enter event loop in main thread
    for {
        select {
        case <-interrupt:
            log.Println("Interrupt, closing")

            if shClient.IsOpen() == false {
                log.Println("ship client not connected, exiting...")
                return
            }

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

    showHelp := getopt.BoolLong("help", 'h', "display this help text")
    showVersion := getopt.BoolLong("version", 'v', "display this help text")
    configFile := getopt.StringLong("config", 'c', "./config.json", "Config file to read", "file")
    pidFile := getopt.StringLong("pid", 'p', "", "Where to write process id", "file")

    getopt.Parse()

    if *showHelp {
        getopt.Usage()
        return
    }

    if *showVersion {
        fmt.Println("v0.0.0")
        return
    }

    // Write PID file
    if len(*pidFile) > 0 {
        log.Printf("Writing pid to: %s", *pidFile)
        err = pid.Save(*pidFile)
        if err != nil {
            log.Println(err)
            return
        }
    }

    // Parse config
    config, err = LoadConfig(*configFile)
    if err != nil {
        log.Println(err)
        return
    }

    // Connect to redis
    err = RedisConnect(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
    if err != nil {
        log.Println("Failed to connect to redis:", err)
        return
    }

    // Connect client and get chain info.
    log.Printf("Get chain info from api at: %s", config.Api)
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

    // Run the application
    run()
}
