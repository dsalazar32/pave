package commands

import (
	"errors"
	"testing"
)

func TestLanguage_Validate(t *testing.T) {
	type test struct {
		lang Language
		want interface{}
	}

	cases := []test{
		{Language{"node", "10.13.0"}, "node 10.13.0"},
		{Language{"node", "8.10.0"}, "node 8.10.0"},
		{Language{"node", "4.3.2"}, "node 4.3.2"},
		{Language{"ruby", "2.1.7"}, "ruby 2.1.7"},
		{Language{"ruby", "2.1.5"}, "ruby 2.1.5"},
		{Language{"node", ""}, "node 10.13.0"},
		{Language{"ruby", ""}, "ruby 2.1.7"},
		{Language{"haskel", "1.2.4"}, errors.New("unsupported language: haskel 1.2.4")},
		{Language{"node", "4.3.0"}, errors.New("unsupported version: node 4.3.0")},
		{Language{"ruby", "zing4.3.0"}, errors.New("error parsing language version: ruby zing4.3.0")},
	}

	for _, c := range cases {
		if err := c.lang.Validate(); err != nil {
			got, want := err.Error(), c.want.(error).Error()
			if got != want {
				t.Fatalf("incorrect error returned: want %v, but got %v", want, got)
			}
		} else {
			got := c.lang.String()
			if got != c.want {
				t.Fatalf("error parsing semver: want %s, but got %s", c.want, got)
			}
		}
	}
}
