package commands

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/coreos/go-semver/semver"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
	"strings"
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

type Support struct {
	SupportedLanguages `yaml:"supported_languages"`
}

type Languages []SupportedLanguage

type SupportedLanguages struct {
	Default   string
	Languages map[string]Languages
}

type SupportedLanguage struct {
	Name      string
	Version   semver.Version
	BaseImage string
}

func (s *Support) Parse(b []byte) error {
	if err := yaml.Unmarshal(b, s); err != nil {
		return fmt.Errorf("error unmarshalling yaml: %s", err)
	}

	// Sort all supported languages by semver
	for _, v := range s.SupportedLanguages.Languages {
		v.Sort()
	}

	return nil
}

func (sl SupportedLanguage) String() string {
	return fmt.Sprintf("%s %s", sl.Name, sl.Version)
}

func (sl SupportedLanguages) Validate(lang string, force bool) (SupportedLanguage, error) {
	lv := make([]string, 2)
	for idx, val := range strings.Split(lang, ":") {
		lv[idx] = val
	}
	lname, lvers := lv[0], lv[1]

	langVers, ok := sl.Languages[lname]
	if !ok {
		return SupportedLanguage{}, fmt.Errorf("unsupported language: %s", lang)
	}
	if len(lvers) == 0 {
		return langVers.Latest(), nil
	}
	for _, lvs := range langVers {
		v, err := semver.NewVersion(lvers)
		if err != nil {
			return SupportedLanguage{}, fmt.Errorf("error parsing language version: %s", lang)
		}
		if lvs.Version.Equal(*v) {
			if !force && !v.Equal(langVers.Latest().Version) {
				return SupportedLanguage{}, fmt.Errorf("error not the latest version [%s]."+
					" To use the desired version pass the `-f` flag", lvers)

			}
			return lvs, nil
		}
	}
	return SupportedLanguage{}, fmt.Errorf("unsupported version: %s", lang)
}

func (l Languages) Len() int {
	return len(l)
}

func (l Languages) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Languages) Less(i, j int) bool {
	return l[i].Version.LessThan(l[j].Version)
}

func (l Languages) Sort() {
	sort.Sort(l)
}

func (l Languages) Latest() SupportedLanguage {
	return l[len(l)-1]
}
