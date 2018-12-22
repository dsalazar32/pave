package config

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestSupportedLanguages_Include(t *testing.T) {

	type tt struct {
		lang  string
		force bool
		want  interface{}
	}

	tc := []tt{
		{"node:10.13.0", false, "node:10.13.0"},
		{"node:10.5.0", true, "node:10.5.0"},
		{"node:8.4.0", true, "node:8.4.0"},
		{"ruby:2.1.7", false, "ruby:2.1.7"},
		{"node", false, "node:10.13.0"},
		{"ruby", false, "ruby:2.1.7"},
		{"haskel:1.2.4", false, errors.New("unsupported language: haskel:1.2.4")},
		{"node:4.3.0", false, errors.New("unsupported version: node:4.3.0")},
		{"ruby:zing4.3.0", true, errors.New("error parsing language version: ruby:zing4.3.0")},
		{"node:10.5.0", false, errors.New("error not the latest version [10.5.0]. To use the" +
			" desired version pass the `-f` flag")},
	}

	s, err := LoadSupportFile(filepath.Join("..", "fixtures", "support.yml"))
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range tc {
		lang, err := s.SupportedLanguages.Include(c.lang, c.force)
		if err != nil {
			got, want := err.Error(), c.want.(error).Error()
			if got != want {
				t.Fatalf("incorrect error returned: want %v, but got %v", want, got)
			}
		} else {
			got, want := lang, c.want
			if got != c.want {
				t.Fatalf("error parsing semver: want %s, but got %s", want, got)
			}
		}
	}
}
