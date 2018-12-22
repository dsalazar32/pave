package config

import (
	"fmt"
)

type Outfile struct {
	Outfile string
	Content string
}

type Outfiles []Outfile

type Platform map[string]Outfiles

func (p Platform) Get(n string) (*Outfiles, error) {
	outs, ok := p[n]
	if !ok {
		return nil, fmt.Errorf("unsupported platform: %s", n)
	}
	return &outs, nil
}

type SupportFiles map[string]Platform

func (sf SupportFiles) For(l string) (*Platform, error) {
	p, ok := sf[l]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", l)
	}
	return &p, nil
}
