package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

type InitCommand struct {
	Meta
}

var (
	fs = flag.NewFlagSet("Init", flag.ContinueOnError)

	project = fs.String("p", "",
		"Set the `project` name. Will default to the name of the project directory.")

	projectLang = fs.String("l", "",
		"Set the project language `LANG[:VERSION]`.\n"+
			"If LANG is provided the most recently supported version is defaulted.\n"+
			"If LANG[:VERSION] is other than the most recently supported version, it is required that you set -f flag as well.")

	force = fs.Bool("f", false,
		"Set the force flag to allow validation to use previous supported versions of the supported languages.")

	dockerEnabled = fs.Bool("with-docker", false,
		"Generate Dockerfile for the supported language.")

	terraformEnabled = fs.Bool("with-terraform", false,
		"Enables terraform support.\n"+
			"When enabled tinker will generate an initial policy that will create ECS resources for the project to land on.\n"+
			"There will be support for adding other project dependencies further down the line (ex. s3, rds, dynamo, etc...).")

	printLanguages = fs.Bool("list-languages", false,
		"Print supported languages and exit")

	printAllVersions = fs.Bool("all-versions", false,
		"Print supported languages with all versions and exit")
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

			support := Support{}
			if err := support.Parse("./support.yml"); err != nil {
				c.Ui.Error(err.Error())
				return 1
			}
			fmt.Printf("%#v\n", support.SupportedLanguages)

			// Print Languages Version(s)
			if *printLanguages {
				vHeader := "version"
				if *printAllVersions {
					vHeader = vHeader + "s"
				}
				rows := [][]string{
					{"language", vHeader},
					{"--------", strings.Repeat("-", len(vHeader))},
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
				return 0
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

			// TODO: See about inferring language via .<lang>-version files if present.
			// Configure Project Language
			// - Set via `-l` flag
			// - Set via `ask` if `-l` flag is not used. Defaults to `DefaultLanguage` const [node].
			if len(*projectLang) == 0 {
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

			// Enable Dockerfile support
			if !*dockerEnabled {
				b, err := c.Ui.Ask("Enable docker support? [false]")
				if err != nil {
					fmt.Println(err)
					return 1
				}
				if len(b) > 0 {
					*dockerEnabled, err = strconv.ParseBool(b)
					if err != nil {
						fmt.Println(err)
						return 1
					}
				}
			}
			cfg.Dockerfile = *dockerEnabled

			// Enable Terraform support
			if !*terraformEnabled {
				b, err := c.Ui.Ask("Enable terraform support? [false]")
				if err != nil {
					fmt.Println(err)
					return 1
				}
				if len(b) > 0 {
					*terraformEnabled, err = strconv.ParseBool(b)
					if err != nil {
						fmt.Println(err)
						return 1
					}
				}
			}
			cfg.Terraform = *terraformEnabled

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
