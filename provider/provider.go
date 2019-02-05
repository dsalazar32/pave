package provider

import (
	"fmt"
	"strings"
)

type Providers int

const (
	FSYS Providers = iota
	AWS
	GCP
)

type Provider interface {
	Read() (string, error)
}

var pmap = map[string]Providers{
	"s3": AWS,
	"gs": GCP,
}

type ProviderSpec struct {
	New         func(infile string) (Provider, error)
	description string
}

var Constructors = map[Providers]*ProviderSpec{}

func ProviderLookup(s string) (*ProviderSpec, error) {
	p := strings.Index(s, "://")
	if p == -1 {
		return Constructors[FSYS], nil
	} else {
		if provider, ok := pmap[s[:p]]; ok {
			return Constructors[provider], nil
		} else {
			return nil, fmt.Errorf("Error looking up storage provider for %s:", provider)
		}
	}
}
