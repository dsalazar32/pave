package config

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
	"strings"
)

type SupportedLanguages struct {
	Default   string
	Languages Languages
}

func (s SupportedLanguages) Include(l string, force bool) (string, error) {
	langVers := strings.Split(l, ":")
	pinned := len(langVers) > 1

	versions, ok := s.Languages[langVers[0]]
	if !ok {
		return "", fmt.Errorf("unsupported language: %s", l)
	}
	latestVersion := versions.Latest()

	if !pinned {
		return latestVersion.String(), nil
	}

	sver, err := semver.NewVersion(langVers[1])
	if err != nil {
		return "", fmt.Errorf("error parsing language version: %s", l)
	}

	if sver.Equal(latestVersion.Version) {
		return latestVersion.String(), nil
	}

	for _, v := range versions {
		if sver.Equal(v.Version) {
			if !force {
				return "", fmt.Errorf("error not the latest version [%s]."+
					" To use the desired version pass the `-f` flag", sver)
			} else {
				return v.String(), nil
			}
		}
	}

	return "", fmt.Errorf("unsupported version: %s", l)
}

func (s SupportedLanguages) Show(all bool) {
	vHeader := "version"
	if all {
		vHeader = vHeader + "s"
	}

	rows := [][]string{
		{"language", vHeader},
	}

	for l, vers := range s.Languages {
		var version string
		if all {
			var vs []string
			for _, v := range vers {
				vs = append(vs, v.Version.String())
			}
			version = strings.Join(vs, " ")
		} else {
			version = vers.Latest().Version.String()
		}
		rows = append(rows, []string{l, version})
	}
	printTable(rows)
}

type Languages map[string]Versions

type Versions []Version

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v Versions) Less(i, j int) bool {
	return v[i].Version.LessThan(v[j].Version)
}

func (v Versions) Latest() Version {
	return v[len(v)-1]
}

type Version struct {
	Name      string
	Version   semver.Version
	BaseImage string
}

func (v Version) String() string {
	return fmt.Sprintf("%s:%s", v.Name, v.Version)
}
