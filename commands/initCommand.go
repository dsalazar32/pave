package commands

import (
	"flag"
	"fmt"
	"github.com/dsalazar32/tinker/config"
	"os"
	"path/filepath"
	"strconv"
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
			"When enabled pave will generate an initial policy that will create ECS resources for the project to land on.\n"+
			"There will be support for adding other project dependencies further down the line (ex. s3, rds, dynamo, etc...).")

	printLanguages = fs.Bool("list-languages", false,
		"Print supported languages and exit")

	printAllVersions = fs.Bool("all-versions", false,
		"Print supported languages with all versions and exit")
)

// TODO: Check for .pave directory and files
// If not create all required files and directories
// - .pave
// - .pave_infra
// - Dockerfile
func (c InitCommand) Run(args []string) int {

	if c.Config.IsValid() {
		c.Ui.Info("Pave seems to be initialized in this project.\n" +
			"To re-initialize the project run pave init with the '-force' flag or remove the .pave.yml.")
		return 0
	}

	if err := fs.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("error parsing arguments: %s", err))
		return 1
	}

	sf, err := config.LoadSupportFile("./support.yml")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error loading support file: %s", err))
		return 1
	}

	if *printLanguages {
		sf.SupportedLanguages.Show(*printAllVersions)
		return 0
	}

	c.Config.Pave = &config.Pave{}
	p := c.Config.Pave

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
	p.ProjectName = *project

	// TODO: See about inferring language via .<lang>-version files if present.
	// Configure Project Language
	// - Set via `-l` flag
	// - Set via `ask` if `-l` flag is not used. Defaults to `DefaultLanguage` const [node].
	if len(*projectLang) == 0 {
		defaultLang := sf.SupportedLanguages.Default
		*projectLang, err = c.Ui.Ask(fmt.Sprintf("Project language? [%s]", defaultLang))
		if err != nil {
			fmt.Println(err)
			return 1
		}
		if len(*projectLang) == 0 {
			*projectLang = defaultLang
		}
	}

	lang, err := sf.SupportedLanguages.Include(*projectLang, *force)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	p.ProjectLang = lang

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
	p.DockerEnabled = *dockerEnabled

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
	p.TerraformEnabled = *terraformEnabled

	c.Config.Pave = p

	if err := c.Config.WriteFile(); err != nil {
		c.Ui.Error(fmt.Sprintf("error writing config file: %s\n", err))
		return 1
	}

	return 0
}

func (c InitCommand) Help() string {
	return "go init"
}

func (c InitCommand) Synopsis() string {
	return "go init"
}
