package config

import (
	"os"
	"strings"
	"text/template"
)

type TemplatePackage struct {
	Config  Config
	Support Support
}

var funcMap = template.FuncMap{
	"lookup": func(m map[string]interface{}, k string) interface{} {
		return nil
	},
	"split": func(s, d string) []string {
		return strings.Split(s, d)
	},
}

func ParseTemplate(s string, tp TemplatePackage) (string, error) {
	tmpl, err := template.New(tp.Config.Pave.ProjectName).Funcs(funcMap).Parse(s)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(os.Stdout, tp)
	return "", nil
}
