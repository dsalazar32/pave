package envparser

import (
	"bufio"
	"regexp"
	"strings"
)

type Env struct {
	Name  string
	Value string
}

type Envs []Env

func ParseEnvString(s string) Envs {
	valid := regexp.MustCompile(`^(export )?\w+=.+`)
	envs := Envs{}
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		t := scanner.Text()
		if !valid.MatchString(t) {
			continue
		}

		re := regexp.MustCompile(`\w+=.+`)
		kv := re.FindAllString(t, -1)[0]
		idx := strings.Index(kv, "=")
		k, v := strings.TrimSpace(t[:idx]), t[idx+1:]
		envs = append(envs, Env{k, v})
	}
	return envs
}
