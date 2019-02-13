package provider

import (
	"io/ioutil"
	"os"
)

type FileSystem struct {
	filePath string
}

func init() {
	Constructors[FSYS] = &ProviderSpec{
		New:         NewFileSystem,
		description: "FileSystem provider is use to interact with objects on the OS file system.",
	}
}

func (p FileSystem) Read() (string, error) {
	if _, err := os.Stat(p.filePath); err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(p.filePath)
	return string(b), err
}

func (p FileSystem) Write(s string) error {
	return nil
}

func NewFileSystem(infile string) (Provider, error) {
	return &FileSystem{infile}, nil
}
