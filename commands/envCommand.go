package commands

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dsalazar32/pave/helper/envparser"
	"github.com/dsalazar32/pave/provider"
	"io"
	"os"
	"reflect"
)

type EnvCommand struct {
	Meta
}

func (c EnvCommand) Run(args []string) int {

	var (
		fs = flag.NewFlagSet("Env", flag.ContinueOnError)
		in = fs.String("in", "",
			"Environmental variables properties file.")
	)

	if err := fs.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("error parsing arguments: %s", err))
		return 1
	}

	// TODO: Need to work on documentation and appropriate error verbiage
	// STDIN or -in flag
	var (
		envs    envparser.Envs
		envsErr error
	)
	stdInf, _ := os.Stdin.Stat()
	if (stdInf.Mode()&os.ModeCharDevice) == os.ModeCharDevice && *in == "" {
		c.Ui.Error("Either pipe the property file in or use the `-in` flag.")
		return 1
	} else if (stdInf.Mode()&os.ModeCharDevice) != os.ModeCharDevice && *in != "" {
		c.Ui.Error("STDIN and the `-in` flag cannot be used at the same time.")
		return 1
	} else if stdInf.Size() > 0 {
		b := &bytes.Buffer{}
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil && err == io.EOF {
				break
			}
			b.WriteString(input)
		}
		envs = envparser.ParseEnvString(b.String())
	} else if *in != "" {
		envs, envsErr = loadByProvider(*in)
		if envsErr != nil {
			c.Ui.Error(envsErr.Error())
			return 1
		}
	}
	spew.Dump(envs)

	return 0
}

func (c EnvCommand) Help() string {
	panic("implement me")
}

func (c EnvCommand) Synopsis() string {
	panic("implement me")
}

func loadByProvider(infile string) (envparser.Envs, error) {
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

	var parsed envparser.Envs

	t := reflect.TypeOf(fsys)
	switch t.Elem().Name() {
	case "Chef":
		parsed = envparser.ParseEnvChefJson(file)
	default:
		parsed = envparser.ParseEnvString(file)
	}

	return parsed, nil
}
