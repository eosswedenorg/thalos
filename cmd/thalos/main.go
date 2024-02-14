package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var VersionString string = "dev"

func main() {
	cli.AppHelpTemplate = `Usage: {{.HelpName}} [options]

   {{range .VisibleFlags}}{{.}}
   {{end}}`

	cli.HelpFlag = &cli.BoolFlag{
		Name:               "help",
		Aliases:            []string{"h"},
		Usage:              "display this help text",
		DisableDefaultText: true,
	}

	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Printf("Version %s\n", ctx.App.Version)
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Aliases:            []string{"v"},
		Usage:              "display the version",
		DisableDefaultText: true,
	}

	app := &cli.App{
		Version:                VersionString,
		Args:                   true,
		UseShortOptionHandling: true,
		HideHelpCommand:        true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:      "config",
				Aliases:   []string{"c"},
				Value:     "./config.yml",
				Usage:     "Config `file` to read",
				TakesFile: true,
			},
			&cli.StringFlag{
				Name:    "level",
				Aliases: []string{"L"},
				Usage:   "Log level to use",
				Value:   "info",
			},
			&cli.PathFlag{
				Name:      "log",
				Aliases:   []string{"l"},
				Usage:     "Path to log `file`",
				TakesFile: true,
			},
			&cli.BoolFlag{
				Name:               "n",
				Usage:              "Force the application to take start block from config/api",
				DisableDefaultText: true,
			},
			&cli.StringFlag{
				Name:      "pid",
				Aliases:   []string{"p"},
				Usage:     "`file` to save process id to",
				TakesFile: true,
			},
		},
		Action: serverCmd,
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("Application error")
	}
}
