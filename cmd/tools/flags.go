package main

import (
	"github.com/urfave/cli/v2"
)

var redisPrefixFlag = &cli.StringFlag{
	Name:  "prefix",
	Value: "ship",
}

var redisUrlFlag = &cli.StringFlag{
	Name:  "redis-url",
	Value: "127.0.0.1:6379",
	Usage: "host:port to the redis server",
}

var redisUserFlag = &cli.StringFlag{
	Name:  "redis-user",
	Usage: "User to use when authenticating to the server",
}

var redisPasswordFlag = &cli.StringFlag{
	Name:  "redis-pw",
	Usage: "Password to use when authenticating to the server",
}

var redisDbFlag = &cli.IntFlag{
	Name:  "redis-db",
	Value: 0,
	Usage: "What redis database we should connect to.",
}

var chainIdFlag = &cli.StringFlag{
	Name:  "chain_id",
	Value: "1064487b3cd1a897ce03ae5b6a865651747e2e152090f99c1d19d44e01aea5a4",
}
