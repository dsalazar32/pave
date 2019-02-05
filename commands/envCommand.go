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
	var envs envparser.Envs
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
		p, err := provider.ProviderLookup(*in)
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
		fsys, err := p.New(*in)
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
		file, err := fsys.Read()
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
		envs = envparser.ParseEnvString(file)
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
