package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Params struct {
	ConfigPath string
	ListenOn   string
	Debug      bool
}

func NewCommand() *cobra.Command {
	params := Params{}

	command := &cobra.Command{
		Use:   "webhook-conversion-service",
		Short: "Listen for incoming webhooks and perform conversions",
		Run: func(command *cobra.Command, args []string) {
			if params.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			}

			if err := Main(params.ConfigPath, params.ListenOn, HttpServer{}); err != nil {
				logrus.Fatal(err)
			}
			return
		},
	}

	command.Flags().StringVarP(&params.ConfigPath, "config", "c", "/etc/webhook-conversion-service/config.yaml", "Path to configuration file")
	command.Flags().StringVarP(&params.ListenOn, "listen", "l", ":8080", "Address and port to listen on")
	command.Flags().BoolVarP(&params.Debug, "debug", "d", false, "Set logging level to debug")

	return command
}
