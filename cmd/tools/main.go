package main

import (
	"time"

	_ "github.com/eosswedenorg/thalos/internal/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var VersionString string = "dev"

var rootCmd *cobra.Command

func init() {
	redisFlags := pflag.FlagSet{}
	redisFlags.String("redis-url", "127.0.0.1:6379", "host:port to the redis server")
	redisFlags.String("redis-user", "", "User to use when authenticating to the server")
	redisFlags.String("redis-pw", "", "Password to use when authenticating to the server")
	redisFlags.Int("redis-db", 0, "What redis database we should connect to.")
	redisFlags.String("prefix", "ship", "redis prefix")
	redisFlags.String("chain_id", "1064487b3cd1a897ce03ae5b6a865651747e2e152090f99c1d19d44e01aea5a4", "chain id")

	rootCmd = &cobra.Command{
		Use:   "thalos-tools",
		Short: "Collection of tools for dealing with the thalos application",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		Version: VersionString,
	}

	benchCmd.Flags().AddFlagSet(&redisFlags)
	benchCmd.Flags().DurationP("interval", "i", time.Minute, "How often the benchmark results should be displayed.")

	validateCmd.Flags().AddFlagSet(&redisFlags)

	MockPublisherCmd.Flags().AddFlagSet(&redisFlags)
	MockPublisherCmd.Flags().String("codec", "json", "codec to use")

	RedisACLCmd.Flags().String("default-pw", "", "Password to use for the default account, if not provided a random one will be generated")
	RedisACLCmd.Flags().String("client", "thalos-client", "Thalos client account name")
	RedisACLCmd.Flags().String("client-pw", "", "Password to use for the thalos client account, if not provided a random one will be generated")
	RedisACLCmd.Flags().String("server", "thalos", "Thalos account name")
	RedisACLCmd.Flags().String("server-pw", "", "Password to use for the thalos server account, if not provided a random one will be generated")
	RedisACLCmd.Flags().String("prefix", "ship", "Redis key prefix")
	RedisACLCmd.Flags().Bool("cleartext", false, "If passwords should be hashed or left in cleartext.")
	RedisACLCmd.Flags().String("file", "", "Where the config should be written to (default: standard out)")

	rootCmd.AddCommand(
		validateCmd,
		benchCmd,
		RedisACLCmd,
		MockPublisherCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
