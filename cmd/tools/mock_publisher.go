package main

import (
	"context"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	api_redis "github.com/eosswedenorg/thalos/api/redis"
	redis_driver "github.com/eosswedenorg/thalos/app/driver/redis"

	"github.com/redis/go-redis/v9"
)

var MockPublisherCmd = &cli.Command{
	Name:  "mock_publisher",
	Usage: "Run a publisher that mocks messages to a redis server. tries to send as many messages as possible",
	Flags: []cli.Flag{
		redisUrlFlag,
		redisUserFlag,
		redisPasswordFlag,
		redisDbFlag,
		redisPrefixFlag,
		chainIdFlag,
		&cli.StringFlag{
			Name:  "codec",
			Value: "json",
		},
	},
	Action: func(ctx *cli.Context) error {
		// Create redis client
		rdb := redis.NewClient(&redis.Options{
			Addr:     ctx.String("redis-url"),
			Username: ctx.String("redis-user"),
			Password: ctx.String("redis-pw"),
			DB:       ctx.Int("redis-db"),
		})

		codec, err := message.GetCodec(ctx.String("codec"))
		if err != nil {
			return err
		}

		ns := api_redis.Namespace{
			Prefix:  ctx.String("redis-prefix"),
			ChainID: ctx.String("chain_id"),
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
			return err
		}
		channel := api.ActionChannel{}.Channel()

		for {
			_ = publisher.Write(channel, payload)
			publisher.Flush()
		}
	},
}
