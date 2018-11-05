package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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
	var (
		prj string
	)

	cfg := Config{}
	if err := cfg.Load(); err != nil {

		// Assumption here is that the project has never been initialized.
		// Interrogate to gather proper information for initialization.
		if os.IsNotExist(err) {
			f := flag.NewFlagSet("Init", flag.ContinueOnError)

			// Set project name
			// - Set via argument flag.
			// - Ask and infer the default value being that of the directory
			// that the project is in.
			f.StringVar(&prj, "p", "", "project name")

			// Set language that will be used in the project.
			// - Set via argument flag.
			// - Ask and infer from .<lang>-version file.
			m, _ := filepath.Glob(".*-versions")
			if len(m) == 0 {

			} else {

			}

			if err := f.Parse(args); err != nil {
				fmt.Printf("Error parsing arguments: %s", err)
				return 1
			}

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
