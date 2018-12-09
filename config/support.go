package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sort"
)

type Support struct {
	SupportedLanguages `yaml:"supported_languages"`
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
