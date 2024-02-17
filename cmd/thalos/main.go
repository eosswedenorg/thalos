package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var VersionString string = "dev"

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{
		Use: "thalos-server",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		Version: VersionString,
		Run:     serverCmd,
	}

	rootCmd.SetHelpTemplate(
		`{{ .Use | trimTrailingWhitespaces}} v{{.Version}}

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
`)
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "v%s" .Version}}` + "\n")

	flags := pflag.FlagSet{}
	flags.StringP("config", "c", "./config.yml", "Config file to read")
	flags.StringP("level", "L", "info", "Log level to use")
	flags.StringP("log", "l", "", "Path to log file (default: print to stdout/stderr)")
	flags.StringP("pid", "p", "", "Where to write process id")
	flags.BoolP("no-state-cache", "n", false, "Force the application to take start block from config/api")

	rootCmd.PersistentFlags().AddFlagSet(&flags)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
