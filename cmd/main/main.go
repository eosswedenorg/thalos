package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"

	"eosio-ship-trace-reader/abi"
	"eosio-ship-trace-reader/app"
	"eosio-ship-trace-reader/config"
	"eosio-ship-trace-reader/transport/redis_common"
	"eosio-ship-trace-reader/transport/redis_pubsub"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"

	eos "github.com/eoscanada/eos-go"
	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
	"github.com/eosswedenorg-go/pid"
	"github.com/pborman/getopt/v2"
)

// ---------------------------
//  Global variables
// ---------------------------

var conf *config.Config

var shClient *shipclient.Client

// Reader states
const (
	RS_CONNECT = 1
	RS_READ    = 2
)

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
					if err := notify.Send(context.Background(), conf.Name, msg); err != nil {
						log.WithError(err).Error("Failed to send notification")
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
				if shErr, ok := err.(shipclient.ClientError); ok {

					// Bail out if socket is closed
					if shErr.Type == shipclient.ErrSockClosed {
						log.Info("Socket closed, Exiting")
						return
					}

					// Reconnect
					if shErr.Type == shipclient.ErrSockRead || shErr.Type == shipclient.ErrNotConnected {
						state = RS_CONNECT
					}
				}

				log.WithError(err).Error("Failed to read from ship")
			}
		}
	}
}

func run() {
	// Spawn reader loop in another thread.
	go readerLoop()

	// Create interrupt channel.
	signals := make(chan os.Signal, 1)

	// Register signal channel to receive signals from the os.
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt
	sig := <-signals
	log.WithField("signal", sig).Info("Signal received")

	if !shClient.IsOpen() {
		log.Info("ship client not connected, exiting...")
		return
	}

	// Cleanly close the connection by sending a close message.
	err := shClient.Shutdown()
	if err != nil {
		log.WithError(err).Info("failed to send close message to ship server")
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
	var chainInfo *eos.InfoResp

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

	// Init telegram notification service
	telegram, err := telegram.New(conf.Telegram.Id)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize telegram")
		return
	}

	telegram.AddReceivers(conf.Telegram.Channel)

	// Register services in notification manager
	notify.UseServices(telegram)

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to redis")
		return
	}

	log.Printf("Get chain info from api at: %s", conf.Api)
	eosClient := eos.New(conf.Api)
	chainInfo, err = eosClient.GetInfo(context.Background())
	if err != nil {
		log.WithError(err).Fatal("Failed to get info")
		return
	}

	if conf.StartBlockNum == shipclient.NULL_BLOCK_NUMBER {
		if conf.IrreversibleOnly {
			conf.StartBlockNum = uint32(chainInfo.LastIrreversibleBlockNum)
		} else {
			conf.StartBlockNum = uint32(chainInfo.HeadBlockNum)
		}
	}

	shClient = shipclient.NewClient(func(c *shipclient.Client) {
		c.StartBlock = conf.StartBlockNum
		c.EndBlock = conf.EndBlockNum
		c.IrreversibleOnly = conf.IrreversibleOnly
	})

	processor := app.SpawnProccessor(
		shClient,
		redis_pubsub.NewPublisher(rdb, redis_common.Namespace{
			Prefix:  conf.Redis.Prefix,
			ChainID: chainInfo.ChainID.String(),
		}),
		abi.NewAbiManager(rdb, eosClient, conf.Redis.CacheID),
	)

	// Run the application
	run()

	// Close the processor properly
	processor.Close()
}
