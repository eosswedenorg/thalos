package main

import (
	"os"

	"github.com/urfave/cli/v2"

	_ "github.com/eosswedenorg/thalos/app/log"
	log "github.com/sirupsen/logrus"
)

var VersionString string = "dev"

func main() {
	app := &cli.App{
		Usage:   "Collection of tools for dealing with the thalos application",
		Version: VersionString,
		Commands: []*cli.Command{
			validateCmd,
			benchCmd,
			RedisACLCmd,
			MockPublisherCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("Application error")
	}
}
