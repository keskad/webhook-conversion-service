package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Params struct {
	ConfigPath string
}

func NewCommand() *cobra.Command {
	params := Params{}

	command := &cobra.Command{
		Use:   "webhook-conversion-service",
		Short: "Listen for incoming webhooks and perform conversions",
		Run: func(command *cobra.Command, args []string) {
			if err := Main(params.ConfigPath); err != nil {
				logrus.Fatal(err)
			}
			return
		},
	}

	command.Flags().StringVarP(&params.ConfigPath, "config", "c", "/etc/webhook-conversion-service/config.yaml", "Path to configuration file")

	return command
}
