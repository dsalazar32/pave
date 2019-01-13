package config

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	ConfigFile = ".pave.yml"
	PaveDir    = ".pave"
	VersionNum = "1"
)

type Config struct {
	Version string
	Pave    *Pave `yaml:"pave"`
}

func InitConfig(b []byte) *Config {
	c := &Config{}

	// Assumption here is that no file exists so initialize
	// an empty config with the supported version.
	if len(b) == 0 {
		c.Version = VersionNum
		return c
	}

	if err := yaml.Unmarshal(b, c); err != nil {
		panic(err)
	}

	return c
}

func LoadFile(fd string) *Config {
	var f []byte

	b := &bytes.Buffer{}

	_, err := os.Stat(fd)
	if err == nil {
		f, err = ioutil.ReadFile(fd)
		if err != nil {
			panic(err)
		}
		b.Write(f)
	}

	return InitConfig(b.Bytes())
}

func (c Config) IsValid() bool {
	return c.Pave != nil
}

func (c Config) ToYaml() ([]byte, error) {
	b := &bytes.Buffer{}

	_, err := b.WriteString("---\n")
	if err != nil {
		return nil, err
	}

	y, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}

	_, err = b.Write(y)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (c Config) WriteFile() error {
	b, err := c.ToYaml()
	if err != nil {
		return err
	}
	if err := os.Mkdir(PaveDir, 0744); err != nil {
		return err
	}
	if err := ioutil.WriteFile(ConfigFile, b, 0644); err != nil {
		return err
	}
	return nil
}

func (c Config) Scaffold(s Support) error {
	sfiles, err := s.SupportFiles.For(c.Pave.ProjectLang)
	if err != nil {
		return err
	}

	if c.Pave.DockerEnabled {
		p, err := sfiles.Get("docker")
		if err != nil {
			return err
		}
		if err := s.WriteFiles(c, *p); err != nil {
			return err
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
