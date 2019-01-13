package config

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
)

const CONFIGFILE = `---
version: "1"
pave:
  project_name: TestPave
  project_lang: node:10.13.0
  docker_enabled: true
  terraform_enabled: true
`

func TestInitConfig(t *testing.T) {
	type tt struct {
		given string
		want  *Config
	}

	tc := []tt{
		{"", &Config{Version: VersionNum}},
		{CONFIGFILE, &Config{VersionNum,
			&Pave{
				ProjectName:      "TestPave",
				ProjectLang:      "node:10.13.0",
				DockerEnabled:    true,
				TerraformEnabled: true,}}},
	}

	for _, c := range tc {
		got := InitConfig([]byte(c.given))
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("error initializing Config\n\nwant:\n%s \n\ngot:\n%s", spew.Sdump(c.want), spew.Sdump(got))
		}
	}
}

func TestConfig_IsValid(t *testing.T) {
	type tt struct {
		given string
		want  bool
	}

	tc := []tt{
		{"", false},
		{CONFIGFILE, true},
	}

	for _, c := range tc {
		cfg := InitConfig([]byte(c.given))
		got := cfg.IsValid()
		if c.want != got {
			t.Errorf("error validating \n\n%s\nwant: %t, got: %t", spew.Sdump(cfg), c.want, got)
		}
	}
}

func TestConfig_ToYaml(t *testing.T) {
	type tt struct {
		given *Config
		want  string
	}

	tc := []tt{
		{&Config{VersionNum,
			&Pave{
				ProjectName:      "TestPave",
				ProjectLang:      "node:10.13.0",
				DockerEnabled:    true,
				TerraformEnabled: true,}},
			CONFIGFILE},
	}

	for _, c := range tc {
		got, err := c.given.ToYaml()
		if err != nil {
			t.Error(err)
			continue
		}

		if string(got) != c.want {
			t.Errorf("error generating yaml:\n\nwant:\n%s\n\ngot:\n%s", c.want, got)
		}
	}
}
