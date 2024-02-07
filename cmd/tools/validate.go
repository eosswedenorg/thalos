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

var validateCmd = &cli.Command{
	Name:  "validate",
	Usage: "Validate a thalos server by following action traces and makes sure that blocks arrive in order.",
	Flags: []cli.Flag{
		redisUrlFlag,
		redisDbFlag,
		prefixFlag,
		chainIdFlag,
	},
	Action: func(ctx *cli.Context) error {
		status_duration := time.Second * 10

		log.WithFields(log.Fields{
			"url":      ctx.String("redis-url"),
			"prefix":   ctx.String("prefix"),
			"chain_id": ctx.String("chain_id"),
			"database": ctx.Int("redis-db"),
		}).Info("Connecting to redis")

		// Create redis client
		rdb := redis.NewClient(&redis.Options{
			Addr: ctx.String("redis-url"),
			DB:   ctx.Int("redis-db"),
		})

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return err
		}

		log.Println("Connected to redis")

		log.Info("Starting validation, following the stream")

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

		block_num := uint32(0)
		timeout := time.Second * 5
		timer := time.NewTicker(timeout)

		go func() {
			for t := range client.Channel() {
				switch msg := t.(type) {
				case error:
					log.WithError(msg).Error("Error when reading stream")
				case message.ActionTrace:
					if block_num > 0 {
						var diff int32 = int32(msg.BlockNum - block_num)
						if diff < 0 || diff > 1 {
							log.WithFields(log.Fields{
								"current_block": block_num,
								"block":         msg.BlockNum,
								"diff":          diff,
							}).Warn("Invalid")
						}
					}
					block_num = msg.BlockNum
					timer.Reset(timeout)
				}
			}
		}()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)

		for {
			select {
			case <-sig:
				fmt.Println("Got interrupt")
				client.Close()
				return nil
			case <-timer.C:
				log.WithField("duration", timeout).
					Warn("Did not get any messages during the defined duration")
			case <-time.After(status_duration):
				log.WithFields(log.Fields{
					"current_block": block_num,
				}).Info("Status")
			}
		}
	},
}
