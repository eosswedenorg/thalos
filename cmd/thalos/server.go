package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v4"
	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
	shipws "github.com/eosswedenorg-go/antelope-ship-client/websocket"
	"github.com/eosswedenorg-go/pid"
	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	_ "github.com/eosswedenorg/thalos/api/message/msgpack"
	api_redis "github.com/eosswedenorg/thalos/api/redis"
	"github.com/eosswedenorg/thalos/internal/abi"
	"github.com/eosswedenorg/thalos/internal/cache"
	"github.com/eosswedenorg/thalos/internal/config"
	driver "github.com/eosswedenorg/thalos/internal/driver/redis"
	. "github.com/eosswedenorg/thalos/internal/log"
	. "github.com/eosswedenorg/thalos/internal/server"
	redis_cache "github.com/go-redis/cache/v9"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/redis/go-redis/v9"
	antelopeapi "github.com/shufflingpixels/antelope-go/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func readerLoop(conf *config.Config, running *bool, shClient *shipclient.Stream, processor *ShipProcessor) {
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
	}).Info("ship client: Connecting with Exponential Backoff")

	connectOp := func() error {
		recon_cnt++

		log.WithFields(log.Fields{
			"url": conf.Ship.Url,
			"try": recon_cnt,
		}).Info("ship client: Connecting")

		if err := shClient.Connect(conf.Ship.Url); err != nil {
			return err
		}

		// Set stream client start block to processors current block
		// Both values should be the same on first connect, but when reconnecting
		// We don't want to start from the beginning
		shClient.StartBlock = processor.GetCurrentBlock()

		return shClient.SendBlocksRequest()
	}

	for *running {

		err := backoff.RetryNotify(connectOp, exp, func(err error, d time.Duration) {
			if recon_cnt >= 3 {
				msg := fmt.Sprintf("Failed to connect to ship at '%s'", conf.Ship.Url)
				if err := notify.Send(context.Background(), conf.Name, msg); err != nil {
					log.WithError(err).Error("Failed to send notification")
				}
				recon_cnt = 0
			}

			log.WithError(err).Error("ship client: Failed to connect")

			log.WithFields(log.Fields{
				"reconn_at": time.Now().Add(d),
				"reconn_in": d,
			}).Info("ship client: Reconnecting in ", d)
		})
		if err != nil {
			log.WithError(err).Error("ship client:Failed to connect")
			return
		}

		recon_cnt = 0
		log.WithFields(log.Fields{
			"start": shClient.StartBlock,
			"end":   shClient.EndBlock,
		}).Info("ship client: Connected to ship")

		if err := shClient.Run(); err != nil {

			if errors.Is(err, shipclient.ErrEndBlockReached) {
				log.Info("ship client: Endblock reached.")
				return
			}

			if shipws.IsCloseError(err) {
				log.Info("ship client: Connection closed normally")
				return
			}

			log.WithError(err).Error("ship client: Failed to read message")
		}
	}
}

func run(conf *config.Config, shClient *shipclient.Stream, processor *ShipProcessor) {
	running := true

	// Spawn reader loop in another thread.
	go readerLoop(conf, &running, shClient, processor)

	// Create interrupt channel.
	signals := make(chan os.Signal, 1)

	// Register signal channel to receive signals from the os.
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt
	sig := <-signals
	log.WithField("signal", sig).Info("Signal received")

	// Cleanly close the connection by sending a close message.
	running = false
	err := shClient.Shutdown()
	if err != nil {
		log.WithError(err).Info("failed to send close message to ship server")
	}
}

func LogLevels() []string {
	list := []string{}
	for _, lvl := range log.AllLevels {
		list = append(list, lvl.String())
	}
	return list
}

func initAbiManager(cfg *config.AbiCache, api *antelopeapi.Client, store cache.Store, chain_id string) *abi.AbiManager {
	cache := cache.NewCache("thalos::cache::abi::"+chain_id, store)
	return abi.NewAbiManager(cfg, cache, api)
}

func stateLoader(conf *config.Config, start_block_flag *pflag.Flag, chainInfo func() *antelopeapi.Info, cache *cache.Cache, current_block_no_cache bool) StateLoader {
	return func(state *State) {
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

				if start_block_flag != nil && start_block_flag.Changed {
					source = "cli"
				} else {
					source = "config"
				}

				state.CurrentBlock = conf.Ship.StartBlockNum
			} else {
				// Otherwise, set from api.
				if conf.Ship.IrreversibleOnly {
					source = "api (LIB)"
					state.CurrentBlock = uint32(chainInfo().LastIrreversableBlockNum)
				} else {
					source = "api (HEAD)"
					state.CurrentBlock = uint32(chainInfo().HeadBlockNum)
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

func stateSaver(cache *cache.Cache) StateSaver {
	return func(state State) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()
		return cache.Set(ctx, "state", state, 0)
	}
}

func GetConfig(flags *pflag.FlagSet) (*config.Config, error) {
	filename, err := flags.GetString("config")
	if err != nil {
		return nil, err
	}

	cfg, err := config.NewBuilder().
		SetConfigFile(filename).
		SetFlags(flags).
		Build()
	if err != nil {
		return nil, err
	}

	logFile, _ := flags.GetString("log")
	if len(logFile) > 0 {
		cfg.Log.Directory = path.Dir(logFile)
		cfg.Log.Filename = path.Base(logFile)
	}

	// If start-block is provided, we should set no-state-cache to true.
	if startBlock := flags.Lookup("start-block"); startBlock != nil && startBlock.Changed {
		if err := flags.Set("no-state-cache", "true"); err != nil {
			return cfg, nil
		}
	}

	cfg.Ship.Blacklist.SetWhitelist(cfg.Ship.BlacklistIsWhitelist)

	return cfg, nil
}

func ConnectRedis(conf *config.RedisConfig) (*redis.Client, error) {
	logEntry := log.WithFields(log.Fields{
		"addr": conf.Addr,
		"db":   conf.DB,
	})

	if len(conf.User) > 0 {
		logEntry = logEntry.WithField("user", conf.User)
	}

	if len(conf.Password) > 0 {
		logEntry = logEntry.WithField("password", strings.Repeat("*", len(conf.Password)))
	}

	logEntry.Info("Connecting to redis")

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Username: conf.User,
		Password: conf.Password,
		DB:       conf.DB,
	})

	return rdb, rdb.Ping(context.Background()).Err()
}

func serverCmd(cmd *cobra.Command, args []string) {
	var err error

	// Parse config
	conf, err := GetConfig(cmd.Flags())
	if err != nil {
		log.WithError(err).Fatal("Failed to read config")
		return
	}

	// Write PID file
	pidFile, err := cmd.Flags().GetString("pid")
	if err != nil {
		log.WithField("file", pidFile).Info("Writing pid to file")
		if err = pid.Save(pidFile); err != nil {
			log.WithError(err).Fatal("Failed to write pid")
			return
		}
	}

	skip_currentblock_cache, _ := cmd.Flags().GetBool("no-state-cache")

	flagLevel, _ := cmd.Flags().GetString("level")
	lvl, err := log.ParseLevel(flagLevel)
	if err == nil {
		log.WithField("value", lvl).Info("Setting log level")
		log.SetLevel(lvl)
	} else {
		log.WithError(err).Warn("Failed to parse level")
	}

	if len(conf.Log.Filename) > 0 {
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

	rdb, err := ConnectRedis(&conf.Redis)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to redis")
		return
	}

	// Setup cache storage
	cacheStore := cache.NewRedisStore(&redis_cache.Options{
		Redis: rdb,
		// Cache 10k keys for 10 minutes.
		LocalCache: redis_cache.NewTinyLFU(10000, 10*time.Minute),
	})

	// Setup general cache
	cache := cache.NewCache("thalos::cache::instance::"+conf.Name, cacheStore)

	antelopeClient := antelopeapi.New(conf.Api)

	shClient := shipclient.NewStream(func(s *shipclient.Stream) {
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

	chainInfo := chainInfoOnce(antelopeClient)

	chain_id := conf.Ship.Chain
	if len(chain_id) < 1 {
		chain_id = chainInfo().ChainID
	}

	processor := SpawnProccessor(
		shClient,
		stateLoader(conf, cmd.Flags().Lookup("start-block"), chainInfo, cache, skip_currentblock_cache),
		stateSaver(cache),
		driver.NewPublisher(context.Background(), rdb, api_redis.Namespace{
			Prefix:  conf.Redis.Prefix,
			ChainID: chain_id,
		}),
		initAbiManager(&conf.AbiCache, antelopeClient, cacheStore, chain_id),
		codec,
	)

	processor.SetBlacklist(conf.Ship.Blacklist)

	// Run the application
	run(conf, shClient, processor)

	// Close the processor properly
	processor.Close()
}
