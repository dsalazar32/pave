package strparser

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"
)

type FuncMap template.FuncMap

type TemplatePackage struct {
	Ns      string
	FuncMap FuncMap
	Data    interface{}
}

func ParseTemplate(s string, tp TemplatePackage, wr io.Writer) error {
	tmpl := template.New(tp.Ns)

	if len(tp.FuncMap) > 0 {
		tmpl = tmpl.Funcs(template.FuncMap(tp.FuncMap))
	}

	tmpl, err := tmpl.Parse(s)
	if err != nil {
		return err
	}

	if err = tmpl.Execute(wr, tp.Data); err != nil {
		fmt.Println(err)
		return errors.New(strings.Trim(strings.Join(strings.Split(err.Error(), ":")[4:], ":"), " "))
	}

	return nil
}
