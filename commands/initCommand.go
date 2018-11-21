package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InitCommand struct {
	Meta
}

type langVers []string

func (l langVers) String() string {
	return ""
}

func (l *langVers) Set(value string) error {
	lang := strings.Split(value, ":")
	if len(lang) > 2 {
		return errors.New("expected [language] or [language:vers]")
	}
	for _, lv := range lang {
		*l = append(*l, lv)
	}
	return nil
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
		p    string
		pl   langVers
		lang Language
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
			f.StringVar(&p, "p", "", "project name")

			// Set language that will be used in the project.
			// - Set via argument flag.
			// - Ask and infer from .<lang>-version file.
			// m, _ := filepath.Glob(".*-versions")
			// if len(m) == 0 {
			// } else {
			// }
			f.Var(&pl, "l", "project language")

			if err := f.Parse(args); err != nil {
				c.Ui.Error(fmt.Sprintf("Error parsing arguments: %s", err))
				return 1
			}

			// Configure Project
			if len(p) == 0 {
				pd, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					return 1
				}
				pd = filepath.Base(pd)
				p, _ = c.Ui.Ask(fmt.Sprintf("Project name? [%s]", pd))
				if len(p) == 0 {
					p = pd
				}
			}
			cfg.ProjectName = p

			// Configure Project Language
			if len(pl) == 0 {
				l, _ := c.Ui.Ask("Project language? [node]")
				if len(l) == 0 {
					l = "node"
				}
				lang = Language{name: l}
			} else {
				lang.Set(pl...)
			}
			if err := lang.Validate(); err != nil {
				c.Ui.Error(err.Error())
				return 1
			}
			cfg.Language = lang.String()

			if err := initializeProject(cfg); err != nil {
				c.Ui.Error(fmt.Sprintf("Something bad happened: %v\n", err))
				return 1
			}
		}
	} else {
		c.Ui.Info("Tinker seems to be initialized in this project.\n" +
			"To re-initialize the project run tinker init with the '-force' flag or remove the .tinker directory.")
	}
	return 0
}

func initializeProject(cfg Config) error {
	fmt.Printf("%v\n", cfg)
	return nil
}

func (c InitCommand) Help() string {
	return "go init"
}

func (c InitCommand) Synopsis() string {
	return "go init"
}
