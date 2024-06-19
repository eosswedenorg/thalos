package main

import (
	_ "github.com/eosswedenorg/thalos/internal/log"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var VersionString string = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:   "thalos-tools",
		Short: "Collection of tools for dealing with the thalos application",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		Version: VersionString,
	}

	rootCmd.AddCommand(
		CreateValidateCmd(),
		CreateBenchCmd(),
		CreateRedisACLCmd(),
		CreateMockPublisherCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
