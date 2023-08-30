package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v4"
	eos "github.com/eoscanada/eos-go"
	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
	"github.com/eosswedenorg-go/pid"
	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	_ "github.com/eosswedenorg/thalos/api/message/msgpack"
	api_redis "github.com/eosswedenorg/thalos/api/redis"
	"github.com/eosswedenorg/thalos/app"
	"github.com/eosswedenorg/thalos/app/abi"
	"github.com/eosswedenorg/thalos/app/config"
	. "github.com/eosswedenorg/thalos/app/log"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/pborman/getopt/v2"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// ---------------------------
//  Global variables
// ---------------------------

var conf *config.Config

var shClient *shipclient.Stream

var running bool = false

var VersionString string = "dev"

var exit chan bool

func readerLoop(processor *app.ShipProcessor) {
	running = true
	recon_cnt := 0

	exp := &backoff.ExponentialBackOff{
		InitialInterval:     time.Second,
		RandomizationFactor: 0.25,
		Multiplier:          2,
		MaxInterval:         10 * time.Minute,
		MaxElapsedTime:      0,
		Stop:                -1,
		Clock:               backoff.SystemClock,
	}
	exp.Reset()

	log.WithFields(log.Fields{
		"initial_interval":     exp.InitialInterval,
		"max_interval":         exp.MaxInterval,
		"randomization_factor": exp.RandomizationFactor,
		"multiplier":           exp.Multiplier,
	}).Info("Connecting with Exponential Backoff")

	connectOp := func() error {
		recon_cnt++

		log.WithFields(log.Fields{
			"url": conf.Ship.Url,
			"try": recon_cnt,
		}).Info("Connecting to ship")

		if err := shClient.Connect(conf.Ship.Url); err != nil {
			return err
		}

		// Set stream client start block to processors current block
		// Both values should be the same on first connect, but when reconnecting
		// We don't want to start from the beginning
		shClient.StartBlock = processor.GetCurrentBlock()

		return shClient.SendBlocksRequest()
	}

	for running {

		err := backoff.RetryNotify(connectOp, exp, func(err error, d time.Duration) {
			if recon_cnt >= 3 {
				msg := fmt.Sprintf("Failed to connect to ship at '%s'", conf.Ship.Url)
				if err := notify.Send(context.Background(), conf.Name, msg); err != nil {
					log.WithError(err).Error("Failed to send notification")
				}
				recon_cnt = 0
			}

			log.WithError(err).Error("Failed to connect to SHIP")

			log.WithFields(log.Fields{
				"reconn_at": time.Now().Add(d),
				"reconn_in": d,
			}).Info("Reconnecting in ", d)
		})
		if err != nil {
			log.WithError(err).Error("Failed to connect to SHIP")
			running = false
			continue
		}

		recon_cnt = 0
		log.WithFields(log.Fields{
			"start": shClient.StartBlock,
			"end":   shClient.EndBlock,
		}).Info("Connected to ship")

		if err := shClient.Run(); err != nil {

			if errors.Is(err, shipclient.ErrEndBlockReached) {
				exit <- true
				log.Info("Endblock reached.")
				break
			}

			log.WithError(err).Error("Failed to read from ship")
		}
	}
}

func run(processor *app.ShipProcessor) {
	// Spawn reader loop in another thread.
	go readerLoop(processor)

	// Create interrupt channel.
	signals := make(chan os.Signal, 1)

	// Register signal channel to receive signals from the os.
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt
	select {
	case sig := <-signals:
		log.WithField("signal", sig).Info("Signal received")

		// Cleanly close the connection by sending a close message.
		err := shClient.Shutdown()
		if err != nil {
			log.WithError(err).Info("failed to send close message to ship server")
		}
	case <-exit:
		// Do nothing, just exit.
	}

	running = false
}

func getChain(def string) string {
	if len(conf.Ship.Chain) > 0 {
		return conf.Ship.Chain
	}
	return def
}

func main() {
	var err error
	var chainInfo *eos.InfoResp

	exit = make(chan bool)

	showHelp := getopt.BoolLong("help", 'h', "display this help text")
	showVersion := getopt.BoolLong("version", 'v', "display this help text")
	configFile := getopt.StringLong("config", 'c', "./config.yml", "Config file to read", "file")
	pidFile := getopt.StringLong("pid", 'p', "", "Where to write process id", "file")
	logFile := getopt.StringLong("log", 'l', "", "Path to log file", "file")

	getopt.Parse()

	if *showHelp {
		getopt.Usage()
		return
	}

	if *showVersion {
		fmt.Println(VersionString)
		return
	}

	// Write PID file
	if len(*pidFile) > 0 {
		log.WithField("file", *pidFile).Info("Writing pid to file")
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

	// If log file is given on the commandline, override config values.
	if len(*logFile) > 0 {
		conf.Log.Directory = path.Dir(*logFile)
		conf.Log.Filename = path.Base(*logFile)
	}

	if len(conf.Log.Filename) > 0 {
		stdWriter, err := NewRotatingFileFromConfig(conf.Log, "info")
		if err != nil {
			log.WithError(err).Fatal("Failed to open info log")
			return
		}
		errWriter, err := NewRotatingFileFromConfig(conf.Log, "error")
		if err != nil {
			log.WithError(err).Fatal("Failed to open error log")
			return
		}

		log.WithFields(log.Fields{
			"maxfilesize":    conf.Log.MaxFileSize,
			"maxage":         conf.Log.MaxTime,
			"directory":      conf.Log.GetDirectory(),
			"info_filename":  stdWriter.GetFilename(),
			"error_filename": errWriter.GetFilename(),
		}).Info("Logging to file")

		log.SetOutput(io.Discard)
		log.AddHook(MakeStdHook(stdWriter))
		log.AddHook(MakeErrorHook(errWriter))
	}

	// Init telegram notification service
	if len(conf.Telegram.Id) > 0 {

		telegram, err := telegram.New(conf.Telegram.Id)
		if err != nil {
			log.WithError(err).Fatal("Failed to initialize telegram")
			return
		}

		telegram.AddReceivers(conf.Telegram.Channel)

		// Register services in notification manager
		notify.UseServices(telegram)
	}

	// Connect to redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Username: conf.Redis.User,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to redis")
		return
	}

	log.WithField("api", conf.Api).Info("Get chain info from api")
	eosClient := eos.New(conf.Api)
	chainInfo, err = eosClient.GetInfo(context.Background())
	if err != nil {
		log.WithError(err).Fatal("Failed to get info")
		return
	}

	if conf.Ship.StartBlockNum == shipclient.NULL_BLOCK_NUMBER {
		if conf.Ship.IrreversibleOnly {
			conf.Ship.StartBlockNum = uint32(chainInfo.LastIrreversibleBlockNum)
		} else {
			conf.Ship.StartBlockNum = uint32(chainInfo.HeadBlockNum)
		}
	}

	shClient = shipclient.NewStream(func(s *shipclient.Stream) {
		s.StartBlock = conf.Ship.StartBlockNum
		s.EndBlock = conf.Ship.EndBlockNum
		s.IrreversibleOnly = conf.Ship.IrreversibleOnly
	})

	// Get codec
	codec, err := message.GetCodec(conf.MessageCodec)
	if err != nil {
		log.WithError(err).Fatal("Failed to load codec")
		return
	}

	chain_id := getChain(chainInfo.ChainID.String())

	processor := app.SpawnProccessor(
		shClient,
		api_redis.NewPublisher(context.Background(), rdb, api_redis.Namespace{
			Prefix:  conf.Redis.Prefix,
			ChainID: chain_id,
		}),
		abi.NewAbiManager(rdb, eosClient, chain_id),
		codec,
	)

	// Run the application
	run(processor)

	// Close the processor properly
	processor.Close()
}
