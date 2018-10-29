package main

import (
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	cli := &cli.CLI{
		Name:     "tinker",
		Version:  "0.0.1",
		Args:     os.Args[1:],
		Commands: initCommands(),
	}
	cli.Run()
}
