package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sort"
	"strings"
)

type Support struct {
	SupportedLanguages `yaml:"supported_languages"`
	SupportFiles       `yaml:"support_files"`
}

func LoadSupportFile(f string) (*Support, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return InitializeSupport(b)
}

func InitializeSupport(b []byte) (*Support, error) {
	s := &Support{}
	if err := yaml.Unmarshal(b, s); err != nil {
		return nil, err
	}
	for _, v := range s.Languages {
		sort.Sort(v)
	}

	return s, nil
}

func (s Support) BaseImageLookup(lv string) string {
	langvers := strings.Split(lv, ":")
	vers := s.SupportedLanguages.Languages[langvers[0]]
	for _, v := range vers {
		if v.Version.String() == langvers[1] {
			return v.BaseImage
		}
	}
	return ""
}
