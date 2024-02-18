package main

import (
	"github.com/eosswedenorg/thalos/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	rootCmd.PersistentFlags().AddFlagSet(config.GetFlags())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
