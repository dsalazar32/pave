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
	Write(string) error
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

func parseCloudStorageUrl(url string) (string, string, error) {
	proto := strings.Index(url, "://")
	urlParts := strings.Split(strings.TrimPrefix(url, url[:proto+3]), "/")
	if len(urlParts) < 2 {
		return "", "", fmt.Errorf("Error parsing Cloud Provider URL: %s", url)
	}
	bucket, key := urlParts[0], strings.Join(urlParts[1:], "/")
	return bucket, key, nil
}
