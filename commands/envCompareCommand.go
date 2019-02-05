package commands

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dsalazar32/pave/helper/envparser"
	"github.com/dsalazar32/pave/provider"
)

type EnvCompareCommand struct {
	Meta
}

func (c EnvCompareCommand) Help() string {
	panic("implement me")
}

func (c EnvCompareCommand) Run(args []string) int {
	if len(args) > 2 || len(args) < 2 {
		c.Ui.Info("Usage pave env <source1> <source2>")
		return 1
	}

	origin, err := parseEnvString(args[0])
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	target, err := parseEnvString(args[1])
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	spew.Dump(origin.Diff(target))

	return 0
}

func (c EnvCompareCommand) Synopsis() string {
	panic("implement me")
}

func parseEnvString(infile string) (envparser.Envs, error) {
	p, err := provider.ProviderLookup(infile)
	if err != nil {
		return nil, err
	}
	fsys, err := p.New(infile)
	if err != nil {
		return nil, err
	}
	file, err := fsys.Read()
	if err != nil {
		return nil, err
	}
	return envparser.ParseEnvString(file), nil
}
