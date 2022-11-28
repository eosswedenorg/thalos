
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
    "eosio-ship-trace-reader/config"
    "eosio-ship-trace-reader/redis"
    "eosio-ship-trace-reader/telegram"
    "eosio-ship-trace-reader/transport"
    eos "github.com/eoscanada/eos-go"
    shipclient "github.com/eosswedenorg-go/eos-ship-client"
)

// ---------------------------
//  Global variables
// ---------------------------

var conf config.Config

var chainInfo *eos.InfoResp

var shClient *shipclient.ShipClient

var eosClient *eos.API
var eosClientCtx = context.Background()

var transporter transport.Driver

// Reader states
const RS_CONNECT = 1
const RS_READ = 2

func readerLoop() {

    state := RS_CONNECT
    var recon_cnt uint  = 0

    for {
        switch state {
        case RS_CONNECT :
            recon_cnt++
            log.Printf("Connecting to ship at: %s (Try %d)", conf.ShipApi, recon_cnt)
            err := shClient.Connect(conf.ShipApi)
            if err != nil {
                log.Println(err)

                if recon_cnt >= 3 {
                    msg := fmt.Sprintf("Failed to connect to ship at '%s'", conf.ShipApi)
                    if err = telegram.Send(msg); err != nil {
                        log.Println(err)
                    }
                    recon_cnt = 0
                }

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
            recon_cnt = 0
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
    conf, err = config.Load(*configFile)
    if err != nil {
        log.Println(err)
        return
    }

    // Init telegram
    err = telegram.Init(conf.Name, conf.Telegram.Id, conf.Telegram.Channel)
    if err != nil {
        log.Println("Failed to initialize telegram", err)
        return
    }

    // Connect to redis
    err = redis.Connect(conf.Redis.Addr, conf.Redis.Password, conf.Redis.DB)
    if err != nil {
        log.Println("Failed to connect to redis:", err)
        return
    }

    // Init Abi cache
    InitAbiCache(conf.Redis.CacheID)

    // Connect client and get chain info.
    log.Printf("Get chain info from api at: %s", conf.Api)
    eosClient = eos.New(conf.Api)
    chainInfo, err = eosClient.GetInfo(eosClientCtx)
    if err != nil {
        log.Println("Failed to get info:", err)
        return
    }

    // Create message queue
    transporter, err = transport.Make(conf.Transport, chainInfo.ChainID.String())
    if err != nil {
        log.Println("Failed to create queue:", err)
        return
    }

    redis.Prefix += chainInfo.ChainID.String() + "."

    if conf.StartBlockNum == config.NULL_BLOCK_NUMBER {

        if conf.IrreversibleOnly {
            conf.StartBlockNum = uint32(chainInfo.LastIrreversibleBlockNum)
        } else {
            conf.StartBlockNum = uint32(chainInfo.HeadBlockNum)
        }
    }

    // Construct ship client
    shClient = shipclient.NewClient(conf.StartBlockNum, conf.EndBlockNum, conf.IrreversibleOnly)
    shClient.BlockHandler = processBlock
    shClient.TraceHandler = processTraces

    // Run the application
    run()
}
