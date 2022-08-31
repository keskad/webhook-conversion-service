package main

import (
	"github.com/keskad/webhook-conversion-service/pkg/app"
	"os"
)

func main() {
	command := app.NewCommand()
	args := os.Args

	if args != nil {
		args = args[1:]
		command.SetArgs(args)
	}

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
