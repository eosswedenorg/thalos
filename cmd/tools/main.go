package main

import (
	"os"

	"github.com/spf13/cobra"

	_ "github.com/eosswedenorg/thalos/app/log"
	log "github.com/sirupsen/logrus"
)

var VersionString string = "dev"

var rootCmd = &cobra.Command{
	Use:     os.Args[0],
	Short:   "Collection of tools for dealing with the thalos application",
	Version: VersionString,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Application error")
	}
}
