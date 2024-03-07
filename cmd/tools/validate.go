package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	api_redis "github.com/eosswedenorg/thalos/api/redis"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a thalos server by following action traces and makes sure that blocks arrive in order.",
	Run: func(cmd *cobra.Command, args []string) {
		status_duration := time.Second * 10

		url, _ := cmd.Flags().GetString("redis-url")
		prefix, _ := cmd.Flags().GetString("prefix")
		chain_id, _ := cmd.Flags().GetString("chain_id")
		db, _ := cmd.Flags().GetInt("redis-db")

		log.WithFields(log.Fields{
			"url":      url,
			"prefix":   prefix,
			"chain_id": chain_id,
			"database": db,
		}).Info("Connecting to redis")

		// Create redis client
		rdb := redis.NewClient(&redis.Options{
			Addr: url,
			DB:   db,
		})

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			log.WithError(err).Fatal("Failed to connect to redis")
			return
		}

		log.Println("Connected to redis")

		log.Info("Starting validation, following the stream")

		sub := api_redis.NewSubscriber(context.Background(), rdb, api_redis.Namespace{
			Prefix:  prefix,
			ChainID: chain_id,
		})

		codec, err := message.GetCodec("json")
		if err != nil {
			log.WithError(err).Fatal("Failed to get codec")
			return
		}

		client := api.NewClient(sub, codec.Decoder)

		// Subscribe to all actions
		if err = client.Subscribe(api.ActionChannel{}.Channel()); err != nil {
			log.WithError(err).Fatal("Failed to subscribe to channels")
			return
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
						diff := int32(msg.BlockNum - block_num)
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
				return
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
