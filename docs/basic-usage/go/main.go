package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/eosswedenorg/thalos/api"
	"github.com/eosswedenorg/thalos/api/message"
	_ "github.com/eosswedenorg/thalos/api/message/json"
	api_redis "github.com/eosswedenorg/thalos/api/redis"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Create redis client
	rdb := redis.NewClient(&redis.Options{})

	sub := api_redis.NewSubscriber(context.Background(), rdb, api_redis.Namespace{
		Prefix:  "ship",
		ChainID: "1064487b3cd1a897ce03ae5b6a865651747e2e152090f99c1d19d44e01aea5a4", // Wax mainnet.
	})

	codec, err := message.GetCodec("json")
	if err != nil {
		fmt.Println("Failed to get json codec")
		return
	}

	client := api.NewClient(sub, codec.Decoder)

	client.OnAction = func(act message.ActionTrace) {
		fmt.Println("ActionTrace")
		fmt.Println(act)
		fmt.Println("---")
	}

	client.OnHeartbeat = func(hb message.HeartBeat) {
		fmt.Println("HeartBeat -- block:", hb.BlockNum, "head:", hb.HeadBlockNum, "lib:", hb.LastIrreversibleBlockNum)
	}

	// Subscribe to some stuffs.
	client.Subscribe(api.ActionChannel{Contract: "eosio"}.Channel())
	client.Subscribe(api.ActionChannel{Name: "mine"}.Channel())
	client.Subscribe(api.HeartbeatChannel)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)

		<-sig
		fmt.Println("Got interrupt")
		client.Close()
	}()

	// Read stuff.
	client.Run()
}
