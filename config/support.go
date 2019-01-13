package config

import (
	"bytes"
	"fmt"
	"github.com/dsalazar32/pave/helper/strparser"
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

func (s Support) WriteFiles(c Config, outs Outfiles) error {
	tmpl := strparser.TemplatePackage{
		Ns: "SupportFiles",
		FuncMap: strparser.FuncMap{
			"baseImage": func(l string) string {
				return s.BaseImageLookup(l)
			},
			"leadTab": func(c string) string {
				return fmt.Sprintf("%s%s", "\t", c)
			},
		},
		Data: c,
	}

	for _, o := range outs {
		b := &bytes.Buffer{}
		if err := strparser.ParseTemplate(o.Content, tmpl, b); err != nil {
			return err
		}
		if err := ioutil.WriteFile(o.Outfile, b.Bytes(), o.Perms); err != nil {
			return err
		}
	}

	return nil
}
