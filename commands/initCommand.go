package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

type InitCommand struct {
	Meta
}

var (
	fs = flag.NewFlagSet("Init", flag.ContinueOnError)

	project = fs.String("p", "", "Set the `project` name. Will default to the name of the"+
		" project directory.")

	projectLang = fs.String("l", "", "Set the project language `LANG[:VERSION]`.\n"+
		"If LANG is provided the most recently supported version is defaulted.\n"+
		"If LANG[:VERSION] is other than the most recently supported version, it is required\n"+
		"that you set -f flag as well.")

	force = fs.Bool("f", false, "Set the force flag to allow validation to use previous supported"+
		" versions of the supported languages.")

	printLanguages = fs.Bool("list-languages", false, "Print supported languages")

	printAllVersions = fs.Bool("show-all-versions", false, "Print supported languages with all versions")
)

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

		// Assumption here is that the project has never been initialized.
		// Interrogate to gather proper information for initialization.
		if os.IsNotExist(err) {
			if err := fs.Parse(args); err != nil {
				c.Ui.Error(fmt.Sprintf("Error parsing arguments: %s", err))
				return 1
			}

			// Configure Project Name
			// - Set via `-p` flag
			// - Set via `ask` if `-p` flag is not used. Defaults to project directory.
			if len(*project) == 0 {
				projectDir, err := os.Getwd()
				if err != nil {
					fmt.Println(err)
					return 1
				}
				projectDir = filepath.Base(projectDir)
				*project, _ = c.Ui.Ask(fmt.Sprintf("Project name? [%s]", projectDir))
				if len(*project) == 0 {
					*project = projectDir
				}
			}
			cfg.ProjectName = *project
			fmt.Println("")

			// TODO: See about inferring language via .<lang>-version files if present.
			// Configure Project Language
			// - Set via `-l` flag
			// - Set via `ask` if `-l` flag is not used. Defaults to `DefaultLanguage` const [node].
			if len(*projectLang) == 0 {

				// Print Languages Version(s)
				if *printLanguages {
					rows := [][]string{
						{"language", "version(s)"},
						{"--------", "-------"},
					}
					for l, vers := range supportedLanuages {
						var version string
						if *printAllVersions {
							var vs []string
							for _, v := range vers {
								vs = append(vs, v.String())
							}
							version = strings.Join(vs, " ")
						} else {
							version = vers[len(vers)-1].String()
						}
						rows = append(rows, []string{l, version})
					}
					printTable(rows)
					fmt.Println("")
				}

				*projectLang, err = c.Ui.Ask(fmt.Sprintf("Project language? [%s]", DefaultLanguage))
				if err != nil {
					fmt.Println(err)
					return 1
				}
				if len(*projectLang) == 0 {
					*projectLang = DefaultLanguage
				}
			}
			lang := Language{}
			lang = lang.Parse(*projectLang)
			if err := lang.Validate(*force); err != nil {
				fmt.Println(err)
				return 1
			}
			cfg.Language = lang.String()
			fmt.Println("")

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

func (c InitCommand) Help() string {
	return "go init"
}

func (c InitCommand) Synopsis() string {
	return "go init"
}

func initializeProject(cfg Config) error {
	fmt.Printf("%+v\n", cfg)
	return nil
}

func printTable(rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
	for _, r := range rows {
		fmt.Fprintln(w, strings.Join(r, "\t"))
	}
	w.Flush()
}
