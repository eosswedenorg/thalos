# Thalos Golang API

## Usage

The api is designed with callback functions that are called when messages arrive

## Example

```go
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

    // Create client
	codec, err := message.GetCodec("json")
	if err != nil {
		fmt.Println("Failed to get json codec")
		return
	}

	client := api.NewClient(sub, codec.Decoder)

    // Subscribe to some channels.
	err = client.Subscribe(
		api.TransactionChannel,
		api.ActionChannel{Contract: "eosio"}.Channel(),
		api.ActionChannel{Name: "mine"}.Channel(),
		api.HeartbeatChannel,
		api.TableDeltaChannel{}.Channel(),
	)

	if err != nil {
		fmt.Println(err)
		return
	}

    // Wait for interrupt in a go routine and close the client.
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, os.Interrupt)

		<-sig
		fmt.Println("Got interrupt")

		client.Close()
	}()

    // Read messages
	for t := range client.Channel() {
		switch msg := t.(type) {
		case error:
			fmt.Println("Error:", msg)
		case message.TransactionTrace:
			fmt.Println("Transaction", msg.BlockNum, msg.ID)
			fmt.Println(msg)
			fmt.Println("---")
		case message.HeartBeat :
			fmt.Println("Heartbeat")
			fmt.Println(msg)
			fmt.Println("---")
		}
	}
}



```

## Message channels and types

There are several types of channels to subscribe to aswell with their respectivly
message types.

NOTE: this is not the same as an go channel. all messages will be posted to the same go channel that can
be accessed by `client.Channel()`

| Channel type         | Message type       | Description                                                                                                          |
| -------------------- | ------------------ | -------------------------------------------------------------------------------------------------------------------- |
| -                    | `error`            | Posted if an error occured on the client. There is no channel to subscribe to. error messages will always be posted. |
| `HeartbeatChannel`   | `HeartBeat`        | Heartbeat message. Used to know if thalos is there or not if messages are not posted frequently on real channels.    |
| `RollbackChannel`    | `RollbackMessage`  | This message is posted if the chain has experienced a microfork.                                                     |
| `TransactionChannel` | `TransactionTrace` | Information about an transaction                                                                                     |
| `ActionChannel`      | `ActionTrace`      | Information about an action                                                                                          |
| `TableDeltaChannel`  | `TableDelta`       | Information about an table change                                                                                    |
