package commands

import (
	"github.com/dsalazar32/tinker/config"
	"github.com/mitchellh/cli"
)

type Meta struct {
	Ui cli.Ui

	Config *config.Config
}
