package commands

import (
	"flag"
	"fmt"
	"github.com/dsalazar32/pave/helper/envparser"
	"github.com/dsalazar32/pave/provider"
)

type EnvPushCommand struct {
	Meta
}

func (c EnvPushCommand) Help() string {
	panic("implement me")
}

func (c EnvPushCommand) Run(args []string) int {

	var (
		fs = flag.NewFlagSet("EnvPush", flag.ContinueOnError)

		nonInteractive = fs.Bool("f", false,
			"The push command is automatically set to interactive. Force flag disables it.")
	)

	if err := fs.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if fs.NArg() > 2 || fs.NArg() < 2 {
		c.Ui.Info("Usage pave env push <source1> <source2>")
		return 1
	}

	a1 := args[0+fs.NFlag()]
	a2 := args[1+fs.NFlag()]

	origin, err := loadByProvider(a1)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	p, err := provider.ProviderLookup(a2)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	f, err := p.New(a2)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	var nEnvs envparser.Envs
	if *nonInteractive {
		nEnvs = origin
	} else {
		fmt.Println("Pushing Environment Variables (interactive)")
		nEnvs = envparser.Envs{}
		for idx, env := range origin {
			name := env.Name
			val := env.Value
			nVal, err := c.Ui.Ask(fmt.Sprintf("%d. %s [%s]", idx+1, name, val))
			if err != nil {
				c.Ui.Error(err.Error())
				return 1
			}
			if len(nVal) != 0 {
				val = nVal
			}
			nEnvs = append(nEnvs, envparser.Env{name, val})
		}
	}

	if err := f.Write(nEnvs.ToString()); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c EnvPushCommand) Synopsis() string {
	panic("implement me")
}
