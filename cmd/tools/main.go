package main

import (
	"os"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var VersionString string = "dev"

var rootCmd = &cobra.Command{
	Use:     os.Args[0],
	Short:   "Collection of tools for dealing with the thalos application",
	Version: VersionString,
}

func init() {
	// Initialize logger
	formatter := log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.0000",
	}

	log.SetFormatter(&formatter)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Application error")
	}
}
