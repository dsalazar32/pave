package storage

import (
	"fmt"
	"strings"
)

type Provider int

const (
	FSYS Provider = iota
	AWS
	GCP
)

var pmap = map[string]Provider{
	"s3": AWS,
	"gs": GCP,
}

type ProviderSpec struct {
	New         func(infile string) (Storage, error)
	description string
}

var Constructors = map[Provider]*ProviderSpec{}

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
