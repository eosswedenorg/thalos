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

var (
	interval time.Duration
	chain_id string

	redis_prefix string
	redis_url    string
	redis_db     int
)

var benchCmd = &cobra.Command{
	Use:     "bench",
	Short:   "Run a benchmark against a thalos node",
	Example: "thalos-tools bench -u 192.168.0.123:6379 --redis-db 1 --chain_id my_id -i 5m",
	Run: func(cmd *cobra.Command, args []string) {
		var counter int = 0

		log.WithFields(log.Fields{
			"url":      redis_url,
			"prefix":   redis_prefix,
			"chain_id": chain_id,
			"database": redis_db,
		}).Info("Connecting to redis")

		// Create redis client
		rdb := redis.NewClient(&redis.Options{
			Addr: redis_url,
			DB:   redis_db,
		})

		status := rdb.Ping(context.Background())

		if status.Err() != nil {
			log.Fatal("cant connect to redis: ", status.Err())
			return
		}

		log.Println("Connected to redis")

		log.WithFields(log.Fields{
			"interval": interval,
		}).Info("Starting benchmark")

		sub := api_redis.NewSubscriber(context.Background(), rdb, api_redis.Namespace{
			Prefix:  redis_prefix,
			ChainID: chain_id,
		})

		codec, err := message.GetCodec("json")
		if err != nil {
			log.Fatal(err)
			return
		}

		client := api.NewClient(sub, codec.Decoder)

		client.OnAction = func(act message.ActionTrace) {
			counter++
		}

		// Subscribe to all actions
		if err = client.Subscribe(api.ActionChannel{}.Channel()); err != nil {
			log.Fatal(err)
			return
		}

		go func() {
			t := time.Now()
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)

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
		}()

		// Read stuff.
		client.Run()
	},
}

func init() {
	benchCmd.Flags().DurationVarP(&interval, "interval", "i", time.Minute, "How often the benchmark results should be displayed.")
	benchCmd.Flags().StringVar(&chain_id, "chain_id", "1064487b3cd1a897ce03ae5b6a865651747e2e152090f99c1d19d44e01aea5a4", "")
	benchCmd.Flags().StringVar(&redis_prefix, "prefix", "ship", "")

	benchCmd.Flags().StringVarP(&redis_url, "redis-url", "u", "127.0.0.1:6379", "host:port to the redis server")
	benchCmd.Flags().IntVar(&redis_db, "redis-db", 0, "What redis database we should connect to.")

	rootCmd.AddCommand(benchCmd)
}
