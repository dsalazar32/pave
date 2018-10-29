package commands

import (
	"github.com/BurntSushi/toml"
	"os"
)

var (
	cPath = ".tinker/config"
)

type Config struct {
	// ProjectName represents the project that tinker was initialized under.
	// Example. If tinker was initialized in a directory by the name
	// of `ProjectFoo` ProjectName will eq ProjectFoo.
	ProjectName string

	// Language will be either inferred if a .(lang)-version file is detected.
	// If the .(lang)-version file is not created then tinker will allow the user
	// to set it based on a list of supported languages. (go, ruby, node, etc...)
	Language string

	// LanguageVersion is the "supported" version of the language.
	LanguageVersion string

	// Enables docker support
	Dockerfile bool

	// Enables terraform support
	Terraform bool
}

func (c *Config) Load() error {
	if _, err := os.Stat(cPath); os.IsExist(err) {
		if _, err := toml.DecodeFile(cPath, &c); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
