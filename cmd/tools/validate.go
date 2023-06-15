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

type Tester struct {
	block_num uint32
	timeout   time.Duration
	timer     *time.Ticker
}

func NewTester(timeout time.Duration) *Tester {
	return &Tester{
		block_num: 0,
		timeout:   timeout,
		timer:     time.NewTicker(timeout),
	}
}

func (t *Tester) OnAction(act message.ActionTrace) {
	if t.block_num > 0 {
		var diff int32 = int32(act.BlockNum - t.block_num)
		if diff < 0 || diff > 1 {
			log.WithFields(log.Fields{
				"current_block": t.block_num,
				"block":         act.BlockNum,
				"diff":          diff,
			}).Warn("Invalid")
		}
	}

	t.block_num = act.BlockNum

	t.timer.Reset(t.timeout)
}

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Run a benchmark against a thalos node",
	Example: "thalos-tools bench -u 192.168.0.123:6379 --redis-db 1 --chain_id my_id -i 5m",
	Run: func(cmd *cobra.Command, args []string) {
		tester := NewTester(time.Second * 5)
		status_duration := time.Second * 10

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

		log.Info("Starting validation, following the stream")

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
		client.OnAction = tester.OnAction

		// Subscribe to all actions
		if err = client.Subscribe(api.ActionChannel{}.Channel()); err != nil {
			log.Fatal(err)
			return
		}

		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)

			for {
				select {
				case <-sig:
					fmt.Println("Got interrupt")
					client.Close()
					return
				case <-tester.timer.C:
					log.WithField("duration", tester.timeout).
						Warn("Did not get any messages during the defined duration")
				case <-time.After(status_duration):
					log.WithFields(log.Fields{
						"current_block": tester.block_num,
					}).Info("Status")
				}
			}
		}()

		// Read stuff.
		client.Run()
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
