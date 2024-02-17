package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
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
	. "github.com/eosswedenorg/thalos/app/cache"
	"github.com/eosswedenorg/thalos/app/config"
	driver "github.com/eosswedenorg/thalos/app/driver/redis"
	. "github.com/eosswedenorg/thalos/app/log"
	redis_cache "github.com/go-redis/cache/v9"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ---------------------------
//  Global variables
// ---------------------------

var conf config.Config

var shClient *shipclient.Stream

var running bool = true

var exit chan bool

var cache *Cache

var cacheStore Store

func readerLoop(processor *app.ShipProcessor) {
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

func LogLevels() []string {
	list := []string{}
	for _, lvl := range log.AllLevels {
		list = append(list, lvl.String())
	}
	return list
}

func initAbiManger(api *eos.API, chain_id string) *abi.AbiManager {
	cache := NewCache("thalos::cache::abi::"+chain_id, cacheStore)
	return abi.NewAbiManager(cache, api)
}

func stateLoader(chainInfo *eos.InfoResp, current_block_no_cache bool) app.StateLoader {
	return func(state *app.State) {
		var source string

		// Load state from cache.
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		err := cache.Get(ctx, "state", &state)
		cancel()

		// on error (cache miss) or if current_block_no_cache is set.
		// set current block from config/api
		if current_block_no_cache || err != nil {
			// Set from config if we have a sane value.
			if conf.Ship.StartBlockNum != shipclient.NULL_BLOCK_NUMBER {
				source = "config"
				state.CurrentBlock = conf.Ship.StartBlockNum
			} else {
				// Otherwise, set from api.
				if conf.Ship.IrreversibleOnly {
					source = "api (LIB)"
					state.CurrentBlock = uint32(chainInfo.LastIrreversibleBlockNum)
				} else {
					source = "api (HEAD)"
					state.CurrentBlock = uint32(chainInfo.HeadBlockNum)
				}
			}
		} else {
			source = "cache"
		}

		log.WithFields(log.Fields{
			"block":  state.CurrentBlock,
			"source": source,
		}).Info("Starting from block")
	}
}

func stateSaver(state app.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	return cache.Set(ctx, "state", state, 0)
}

func ReadConfig(cfg *config.Config, flags *pflag.FlagSet) error {
	filename, err := flags.GetString("config")
	if err != nil {
		return err
	}

	// Read file first.
	if err := cfg.ReadFile(filename); err != nil {
		return err
	}

	// Then override any cli flags
	if err := cfg.ReadCliFlags(flags); err != nil {
		return err
	}

	return nil
}

func serverCmd(cmd *cobra.Command, args []string) {
	var err error
	var chainInfo *eos.InfoResp

	exit = make(chan bool)

	skip_currentblock_cache, _ := cmd.Flags().GetBool("no-state-cache")

	// Write PID file
	pidFile, err := cmd.Flags().GetString("pid")
	if err != nil {
		log.WithField("file", pidFile).Info("Writing pid to file")
		if err = pid.Save(pidFile); err != nil {
			log.WithError(err).Fatal("Failed to write pid")
			return
		}
	}

	// Parse config
	conf = config.New()
	if err = ReadConfig(&conf, cmd.Flags()); err != nil {
		log.WithError(err).Fatal("Failed to read config")
		return
	}

	flagLevel, _ := cmd.Flags().GetString("level")
	lvl, err := log.ParseLevel(flagLevel)
	if err == nil {
		log.WithField("value", lvl).Info("Setting log level")
		log.SetLevel(lvl)
	} else {
		log.WithError(err).Warn("Failed to parse level")
	}

	if len(conf.Log.Filename) > 0 {
		fmt.Println(conf.Log.Filename)

		stdWriter, err := NewRotatingFileFromConfig(conf.Log, "info")
		if err != nil {
			log.WithError(err).Fatal("Failed to set standard log file")
			return
		}
		errWriter, err := NewRotatingFileFromConfig(conf.Log, "error")
		if err != nil {
			log.WithError(err).Fatal("Failed to set error log file")
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
			log.WithError(err).Fatal("Failed to initialize telegram notifier")
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

	// Setup cache storage
	cacheStore = NewRedisStore(&redis_cache.Options{
		Redis: rdb,
		// Cache 10k keys for 10 minutes.
		LocalCache: redis_cache.NewTinyLFU(10000, 10*time.Minute),
	})

	// Setup general cache
	cache = NewCache("thalos::cache::instance::"+conf.Name, cacheStore)

	log.WithField("api", conf.Api).Info("Get chain info from api")
	eosClient := eos.New(conf.Api)
	chainInfo, err = eosClient.GetInfo(context.Background())
	if err != nil {
		log.WithError(err).Fatal("Failed to call eos api")
		return
	}

	shClient = shipclient.NewStream(func(s *shipclient.Stream) {
		s.StartBlock = conf.Ship.StartBlockNum
		s.EndBlock = conf.Ship.EndBlockNum
		s.IrreversibleOnly = conf.Ship.IrreversibleOnly
	})

	// Get codec
	codec, err := message.GetCodec(conf.MessageCodec)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialze codec")
		return
	}

	chain_id := getChain(chainInfo.ChainID.String())

	processor := app.SpawnProccessor(
		shClient,
		stateLoader(chainInfo, skip_currentblock_cache),
		stateSaver,
		driver.NewPublisher(context.Background(), rdb, api_redis.Namespace{
			Prefix:  conf.Redis.Prefix,
			ChainID: chain_id,
		}),
		initAbiManger(eosClient, chain_id),
		codec,
	)

	// Run the application
	run(processor)

	// Close the processor properly
	processor.Close()
}
