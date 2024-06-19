package main

import "github.com/spf13/pflag"

func RedisFlags() *pflag.FlagSet {
	set := pflag.FlagSet{}
	set.String("redis-url", "127.0.0.1:6379", "host:port to the redis server")
	set.String("redis-user", "", "User to use when authenticating to the server")
	set.String("redis-pw", "", "Password to use when authenticating to the server")
	set.Int("redis-db", 0, "What redis database we should connect to.")
	set.String("prefix", "ship", "redis prefix")
	set.String("chain_id", "1064487b3cd1a897ce03ae5b6a865651747e2e152090f99c1d19d44e01aea5a4", "chain id")
	return &set
}
