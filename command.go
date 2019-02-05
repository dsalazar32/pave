package main

import (
	"github.com/dsalazar32/pave/commands"
	"github.com/dsalazar32/pave/config"
	"github.com/mitchellh/cli"
	"os"
)

var meta = commands.Meta{
	Ui:     &cli.BasicUi{Reader: os.Stdin, Writer: os.Stdout},
	Config: config.LoadFile(config.ConfigFile),
}

func initCommands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return commands.InitCommand{Meta: meta}, nil
		},
		"env": func() (cli.Command, error) {
			return commands.EnvCommand{Meta: meta}, nil
		},
		"env compare": func() (cli.Command, error) {
			return commands.EnvCompareCommand{Meta: meta}, nil
		},
	}
}
