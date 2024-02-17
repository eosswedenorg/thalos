package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	api_redis "github.com/eosswedenorg/thalos/api/redis"
	redis_driver "github.com/eosswedenorg/thalos/internal/driver/redis"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var MockPublisherCmd = &cobra.Command{
	Use:   "mock_publisher",
	Short: "Run a publisher that mocks messages to a redis server. tries to send as many messages as possible",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("redis-url")
		user, _ := cmd.Flags().GetString("redis-user")
		pw, _ := cmd.Flags().GetString("redis-pw")
		prefix, _ := cmd.Flags().GetString("prefix")
		chain_id, _ := cmd.Flags().GetString("chain_id")
		db, _ := cmd.Flags().GetInt("redis-db")

		// Create redis client
		rdb := redis.NewClient(&redis.Options{
			Addr:     url,
			Username: user,
			Password: pw,
			DB:       db,
		})

		codecArg, _ := cmd.Flags().GetString("codec")

		codec, err := message.GetCodec(codecArg)
		if err != nil {
			log.WithError(err).Fatal("Failed to get codec")
			return
		}

		log.WithFields(log.Fields{
			"url":      url,
			"prefix":   prefix,
			"chain_id": chain_id,
			"database": db,
		}).Info("Starting mock publisher")

		ns := api_redis.Namespace{
			Prefix:  prefix,
			ChainID: chain_id,
		}
		publisher := redis_driver.NewPublisher(context.Background(), rdb, ns)

		msg := message.ActionTrace{
			TxID:      "401e8a7e5deb18a2a69fc6559f49509a155f4355c85efb69c1c1fab5b60ee532",
			BlockNum:  18237917,
			Timestamp: time.Date(2014, 3, 22, 11, 36, 43, 0, time.UTC),
			Receipt: &message.ActionReceipt{
				Receiver:       "acc1",
				ActDigest:      "4c5c08be612e937564fc526ebb5fadf34ae8c2a571fe9d7cdb3ffcdfc53b0e8d",
				GlobalSequence: 12314,
				RecvSequence:   237187239,
				AuthSequence: []message.AccountAuthSequence{
					{
						Account:  "acc1",
						Sequence: 2732863,
					},
					{
						Account:  "acc2",
						Sequence: 263762,
					},
				},
				CodeSequence: 2327832,
				ABISequence:  12376189,
			},
			Name:     "fake",
			Contract: "fake",
			Receiver: "acc1",
			Data: map[string]interface{}{
				"one": 238771832,
				"two": "str",
			},
			Authorization: []message.PermissionLevel{
				{
					Actor:      "acc1",
					Permission: "active",
				},
				{
					Actor:      "acc2",
					Permission: "owner",
				},
			},
			Except: "err",
			Error:  2,
			Return: []byte{0xbe, 0xef},
		}

		payload, err := codec.Encoder(msg)
		if err != nil {
			log.WithError(err).Fatal("Failed to encode message")
			return
		}
		channel := api.ActionChannel{}.Channel()

		for {
			_ = publisher.Write(channel, payload)
			publisher.Flush()
		}
	},
}
