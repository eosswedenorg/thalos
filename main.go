package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"eosio-ship-trace-reader/config"
	"eosio-ship-trace-reader/redis"
	"eosio-ship-trace-reader/telegram"

	eos "github.com/eoscanada/eos-go"
	shipclient "github.com/eosswedenorg-go/eos-ship-client"
	"github.com/eosswedenorg-go/pid"
	"github.com/pborman/getopt/v2"
)

// ---------------------------
//  Global variables
// ---------------------------

var conf config.Config

var chainInfo *eos.InfoResp

var shClient *shipclient.ShipClient

var (
	eosClient    *eos.API
	eosClientCtx = context.Background()
)

// Reader states
const RS_CONNECT = 1
const RS_READ = 2

func readerLoop() {
	state := RS_CONNECT
	var recon_cnt uint = 0

	for {
		switch state {
		case RS_CONNECT:
			recon_cnt++
			log.Infof("Connecting to ship at: %s (Try %d)", conf.ShipApi, recon_cnt)
			err := shClient.Connect(conf.ShipApi)
			if err != nil {
				log.Println(err)

				if recon_cnt >= 3 {
					msg := fmt.Sprintf("Failed to connect to ship at '%s'", conf.ShipApi)
					if err = telegram.Send(msg); err != nil {
						log.WithError(err).Error("Failed to send to telegram")
					}
					recon_cnt = 0
				}

				log.Info("Trying again in 5 seconds ....")
				time.Sleep(5 * time.Second)
				break
			}

			err = shClient.SendBlocksRequest()
			if err != nil {
				log.Println(err)
				break
			}

			// Connected
			log.Infof("Connected, Start: %d, End: %d", shClient.StartBlock, shClient.EndBlock)
			state = RS_READ
			recon_cnt = 0
		case RS_READ:
			err := shClient.Read()
			if err != nil {
				log.WithError(err).Error("Failed to read from ship")

				if shErr, ok := err.(shipclient.ShipClientError); ok {
					// Reconnect
					if shErr.Type == shipclient.ErrSockRead || shErr.Type == shipclient.ErrNotConnected {
						state = RS_CONNECT
					}
				}
			}
		}
	}
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
			log.Info("Interrupt, closing")

			if !shClient.IsOpen() {
				log.Info("ship client not connected, exiting...")
				return
			}

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := shClient.SendCloseMessage()
			if err != nil {
				log.WithError(err).Info("failed to send close message to ship server")
			}

			select {
			case <-done:
				log.Info("Closed")
			case <-time.After(time.Second * 10):
				log.Info("Timeout")
			}
			return
		case <-done:
			log.Info("Closed")
			return
		}
	}
}

func init() {
	// Initialize logger
	formatter := log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.0000",
	}

	log.SetFormatter(&formatter)
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
		log.Infof("Writing pid to: %s", *pidFile)
		err = pid.Save(*pidFile)
		if err != nil {
			log.WithError(err).Fatal("failed to write pid file")
			return
		}
	}

	// Parse config
	conf, err = config.Load(*configFile)
	if err != nil {
		log.WithError(err).Fatal("failed to read config file")
		return
	}

	// Init telegram
	err = telegram.Init(conf.Name, conf.Telegram.Id, conf.Telegram.Channel)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize telegram")
		return
	}

	// Connect to redis
	err = redis.Connect(conf.Redis.Addr, conf.Redis.Password, conf.Redis.DB)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to redis")
		return
	}

	// Init Abi cache
	InitAbiCache(conf.Redis.CacheID)

	// Connect client and get chain info.
	log.Printf("Get chain info from api at: %s", conf.Api)
	eosClient = eos.New(conf.Api)
	chainInfo, err = eosClient.GetInfo(eosClientCtx)
	if err != nil {
		log.WithError(err).Fatal("Failed to get info")
		return
	}

	redis.SetPrefix(conf.Redis.Prefix, chainInfo.ChainID.String())

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
