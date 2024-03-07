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

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Run a benchmark against a thalos node",
	Run: func(cmd *cobra.Command, args []string) {
		counter := 0
		interval, _ := cmd.Flags().GetDuration("interval")

		url, _ := cmd.Flags().GetString("redis-url")
		user, _ := cmd.Flags().GetString("redis-user")
		pw, _ := cmd.Flags().GetString("redis-pw")
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
			Addr:     url,
			Username: user,
			Password: pw,
			DB:       db,
		})

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			log.WithError(err).Fatal("Failed to connect to redis")
			return
		}

		log.Println("Connected to redis")

		log.WithFields(log.Fields{
			"interval": interval,
		}).Info("Starting benchmark")

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
				return
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
