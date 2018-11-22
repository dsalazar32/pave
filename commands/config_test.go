package commands

import (
	"errors"
	"testing"
)

func TestLanguage_Validate(t *testing.T) {
	type test struct {
		lang  string
		force bool
		want  interface{}
	}

	cases := []test{
		{"node:10.13.0", false, "node 10.13.0"},
		{"node:8.10.0", true, "node 8.10.0"},
		{"node:4.3.2", true, "node 4.3.2"},
		{"ruby:2.1.7", false, "ruby 2.1.7"},
		{"ruby:2.1.5", true, "ruby 2.1.5"},
		{"node", false, "node 10.13.0"},
		{"ruby", false, "ruby 2.1.7"},
		{"haskel:1.2.4", false, errors.New("unsupported language: haskel 1.2.4")},
		{"node:4.3.0", false, errors.New("unsupported version: node 4.3.0")},
		{"ruby:zing4.3.0", true, errors.New("error parsing language version: ruby zing4.3.0")},
		{"node:8.10.0", false, errors.New("error not the latest version [10.13.0]. To use the" +
			" desired version pass the `-f` flag")},
	}

	for _, c := range cases {
		lang := Language{}
		lang = lang.Parse(c.lang)
		if err := lang.Validate(c.force); err != nil {
			got, want := err.Error(), c.want.(error).Error()
			if got != want {
				t.Fatalf("incorrect error returned: want %v, but got %v", want, got)
			}
		} else {
			got := lang.String()
			if got != c.want {
				t.Fatalf("error parsing semver: want %s, but got %s", c.want, got)
			}
		}
	}
}
