package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	api_redis "github.com/eosswedenorg/thalos/api/redis"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var benchCmd = &cli.Command{
	Name:  "bench",
	Usage: "Run a benchmark against a thalos node",
	Flags: []cli.Flag{
		redisUrlFlag,
		redisUserFlag,
		redisPasswordFlag,
		redisDbFlag,
		prefixFlag,
		chainIdFlag,
		&cli.DurationFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Value:   time.Minute,
			Usage:   "How often the benchmark results should be displayed.",
		},
	},
	Action: func(ctx *cli.Context) error {
		var counter int = 0
		interval := ctx.Duration("interval")

		log.WithFields(log.Fields{
			"url":      ctx.String("redis-url"),
			"prefix":   ctx.String("prefix"),
			"chain_id": ctx.String("chain_id"),
			"database": ctx.Int("redis-db"),
		}).Info("Connecting to redis")

		// Create redis client
		rdb := redis.NewClient(&redis.Options{
			Addr:     ctx.String("redis-url"),
			Username: ctx.String("redis-user"),
			Password: ctx.String("redis-pw"),
			DB:       ctx.Int("redis-db"),
		})

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return err
		}

		log.Println("Connected to redis")

		log.WithFields(log.Fields{
			"interval": interval,
		}).Info("Starting benchmark")

		sub := api_redis.NewSubscriber(context.Background(), rdb, api_redis.Namespace{
			Prefix:  ctx.String("prefix"),
			ChainID: ctx.String("chain_id"),
		})

		codec, err := message.GetCodec("json")
		if err != nil {
			return err
		}

		client := api.NewClient(sub, codec.Decoder)

		// Subscribe to all actions
		if err = client.Subscribe(api.ActionChannel{}.Channel()); err != nil {
			return err
		}

		go func() {
			for t := range client.Channel() {
				switch err := t.(type) {
				case message.ActionTrace:
					counter++
				case error:
					log.WithError(err).Error("Error when reading stream")
				}
			}
		}()

		t := time.Now()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)

		// Read stuff.
		for {
			select {
			case <-sig:
				fmt.Println("Got interrupt")
				client.Close()
				return nil
			case now := <-time.After(interval):
				elapsed := now.Sub(t)
				t = now

				log.WithFields(log.Fields{
					"num_messages": counter,
					"elapsed":      elapsed,
					"msg_per_sec":  float64(counter) / elapsed.Seconds(),
					"msg_per_ms":   float64(counter) / float64(elapsed.Milliseconds()),
					"msg_per_min":  float64(counter) / elapsed.Minutes(),
				}).Info("Benchmark results")

				counter = 0
			}
		}
	},
}
