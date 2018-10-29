package commands

import (
	"fmt"
	"os"
)

type InitCommand struct {
	Meta
}

// TODO: Check for .tinker directory and files
// If not create all required files and directories
// - .tinker
// - .tinker/config
// - .tinker/terraform
// - .(lang)-version
// - Dockerfile
func (c InitCommand) Run(args []string) int {
	cfg := Config{}
	if err := cfg.Load(); err != nil {
		if os.IsNotExist(err) {
			if err := initializeProject(cfg); err != nil {
				fmt.Printf("Something bad happened!: %v\n", err)
				return 1
			}
		} else {
			fmt.Printf("Something bad happened!: %v\n", err)
			return 1
		}
	}
	return 0
}

func initializeProject(cfg Config) error {
	panic("implement!")
	return nil
}

func (c InitCommand) Help() string {
	return "go init"
}

func (c InitCommand) Synopsis() string {
	return "go init"
}
