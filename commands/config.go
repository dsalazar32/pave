package commands

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/coreos/go-semver/semver"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const DefaultLanguage = "node"

var (
	cPath = ".tinker/config"

	// TODO: Consider at some point have this supported list be seeded via API call.
	supportedLanuages = map[string]semver.Versions{
		"node": versions("10.13.0", "8.10.0", "4.3.2"),
		"ruby": versions("2.1.7", "2.1.5"),
		"go":   versions("1.11.1"),
	}
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

type Language struct {
	name string
	vers string
}

func (l Language) LatesetVers() *semver.Version {
	vers := supportedLanuages[l.name]
	return vers[len(vers)-1]
}

func (l Language) String() string {
	return fmt.Sprintf("%s %s", l.name, l.vers)
}

func (l *Language) Parse(lang string) Language {
	lv := strings.Split(lang, ":")
	l.name = lv[0]
	if len(lv) > 1 {
		l.vers = lv[1]
	}
	return *l
}

func (l *Language) Validate(f bool) error {
	vers, ok := supportedLanuages[l.name]
	if !ok {
		return fmt.Errorf("unsupported language: %s", l.String())
	}
	if len(l.vers) == 0 {
		l.vers = l.LatesetVers().String()
	} else {
		for _, v := range vers {
			t, err := semver.NewVersion(l.vers)
			if err != nil {
				return fmt.Errorf("error parsing language version: %s", l.String())
			}
			if v.Equal(*t) {
				if !f && !t.Equal(*l.LatesetVers()) {
					return fmt.Errorf("error not the latest version [%s]."+
						" To use the desired version pass the `-f` flag", l.LatesetVers())
				}
				return nil
			}
		}
		return fmt.Errorf("unsupported version: %s", l.String())
	}
	return nil
}

func versions(vers ...string) semver.Versions {
	svers := semver.Versions{}
	for _, v := range vers {
		ver, err := semver.NewVersion(v)
		if err != nil {
			panic(fmt.Sprintf("error parsing versions: %v", err))
		}
		svers = append(svers, ver)
		semver.Sort(svers)
	}
	return svers
}

type Support struct {
	SupportedLanguages struct {
		Default   string
		Languages map[string][]SupportedLanguage
	} `yaml:"supported_languages"`
}

type SupportedLanguage struct {
	Name      string
	Version   semver.Version
	BaseImage string
}

func (s *Support) Parse(dir string) error {
	f, err := ioutil.ReadFile(dir)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(f, s); err != nil {
		return fmt.Errorf("error unmarshalling yaml: %s", err)
	}
	return nil
}
