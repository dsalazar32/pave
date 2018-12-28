package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	ConfigFile = ".pave.yml"
	VersionNum = "1"
)

type Config struct {
	Version string
	Pave    *Pave `yaml:"pave"`
}

func Read(c string) *Config {
	cfg := Config{}

	if _, err := os.Stat(c); err != nil {
		cfg.Version = VersionNum
		return &cfg
	}

	b, err := ioutil.ReadFile(c)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func (c Config) IsValid() bool {
	return c.Pave != nil
}

func (c Config) WriteFile() error {
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("---\n")
	if err != nil {
		return err
	}
	_, err = f.Write(y)
	if err != nil {
		return err
	}
	f.Sync()

	return nil
}

func (c Config) Scaffold(s Support) error {
	sf, err := s.SupportFiles.For(c.Pave.ProjectLang)
	if err != nil {
		return err
	}

	if c.Pave.DockerEnabled {
		p, err := sf.Get("docker")
		if err != nil {
			return err
		}

		for _, f := range *p {
			s, err := ParseTemplate(f.Content, TemplatePackage{c, s})
			if err != nil {
				return err
			}
			fmt.Println(s)
		}
	}
	return nil
}

func printTable(rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 10, ' ', 0)
	for _, r := range rows {
		fmt.Fprintln(w, strings.Join(r, "\t"))
	}
	w.Flush()
}
