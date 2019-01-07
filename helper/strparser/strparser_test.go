package strparser

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

const DOB = "1980-06-03"

func TestParseTemplate(t *testing.T) {
	type tt struct {
		s    string
		pkg  TemplatePackage
		want string
	}

	tcs := []tt{
		{
			`Hello my name is {{.Name}}`,
			templatePackageMock(),
			"Hello my name is David Salazar",
		},
		{
			`My birthday is on {{prettifyDate .Dob}}`,
			templatePackageMock(),
			"My birthday is on June 3, 1980",
		},
		{
			`Which makes me {{.GetAge}}`,
			templatePackageMock(),
			fmt.Sprintf("Which makes me %d", calcAge(ParseDate(DOB))),
		},
		{
			`{{.Interest}}`,
			templatePackageMock(),
			"I like to tinker about with go.",
		},
		{
			`Check out the invalid field {{.Invalid}}`,
			templatePackageMock(),
			"executing \"TemplateMock\" at <.Invalid>: can't evaluate field Invalid in type strparser.MockData",
		},
	}

	for _, tc := range tcs {
		b := &bytes.Buffer{}
		if err := ParseTemplate(tc.s, tc.pkg, b); err != nil {
			if err.Error() != tc.want {
				t.Errorf("error with expected err response for \"%s\" want %s but got: %s", tc.s, tc.want, err)
			}
			continue
		}
		got := b.String()
		if tc.want != got {
			t.Errorf("error parsing template \"%s\" want %s but got: %s", tc.s, tc.want, got)
		}
	}
}

func templatePackageMock() TemplatePackage {
	return TemplatePackage{
		Ns: "TemplateMock",
		FuncMap: FuncMap{
			"prettifyDate": func(d time.Time) string {
				return fmt.Sprintf("%s %d, %d", d.Month(), d.Day(), d.Year())
			},
		},
		Data: MockData{"David Salazar", ParseDate(DOB), "I like to tinker about with go."},
	}
}

type MockData struct {
	Name     string
	Dob      time.Time
	Interest string
}

func (md MockData) GetAge() int {
	return calcAge(md.Dob)
}

func ParseDate(d string) time.Time {
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		panic(err)
	}
	return t
}

func calcAge(d time.Time) int {
	return int((time.Since(d).Hours() / 24) / 365)
}
