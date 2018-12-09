package main

import (
	"github.com/dsalazar32/tinker/commands"
	"github.com/dsalazar32/tinker/config"
	"github.com/mitchellh/cli"
	"os"
)

var meta = commands.Meta{
	Ui:     &cli.BasicUi{Reader: os.Stdin, Writer: os.Stdout},
	Config: config.Read(config.ConfigFile),
}

func initCommands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return commands.InitCommand{Meta: meta}, nil
		},
	}
}
