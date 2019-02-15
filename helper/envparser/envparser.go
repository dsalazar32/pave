package envparser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dsalazar32/pave/provider"
	"regexp"
	"strings"
)

type Env struct {
	Name  string
	Value string
}

type Envs []Env

type EnvMap map[string]string

func ParseEnvChefJson(s string) Envs {
	var dbs provider.DBags
	b := []byte(s)
	if err := json.Unmarshal(b, &dbs); err != nil {
		fmt.Printf("there was an error unmarshalling the json request %v", err)
	}

	rs := bytes.Buffer{}
	for _, item := range dbs.Items {
		envs, ok := item.(map[string]interface{})["envs"]
		if !ok {
			break
		}
		for k, v := range envs.(map[string]interface{}) {
			// TODO See about removing surrounding escaped quote on the values
			rs.WriteString(fmt.Sprintf("%s=%q\n", k, v))
		}
	}

	return ParseEnvString(rs.String())
}

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
		k, v := strings.TrimSpace(kv[:idx]), kv[idx+1:]
		envs = append(envs, Env{k, v})
	}
	return envs
}

func (a Envs) Diff(b Envs) Envs {
	mA := a.ToMap()
	mB := b.ToMap()
	var d Envs
	for k, v := range mA {
		if _, ok := mB[k]; !ok {
			d = append(d, Env{k, v})
		}
	}
	return d
}

func (e Envs) ToMap() EnvMap {
	m := make(EnvMap)
	for _, env := range e {
		m[env.Name] = env.Value
	}
	return m
}

func (e Env) ToString() string {
	return fmt.Sprintf("%s=%s", e.Name, e.Value)
}

func (es Envs) ToString() string {
	b := bytes.Buffer{}
	for i, env := range es {
		b.WriteString(env.ToString())
		if len(es)-1 != i {
			b.WriteString("\n")
		}
	}
	return b.String()
}
