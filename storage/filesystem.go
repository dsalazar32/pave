package storage

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

func (s FileSystem) Read() (string, error) {
	if _, err := os.Stat(s.filePath); err != nil {
		return "", err
	}
	b, err := ioutil.ReadFile(s.filePath)
	return string(b), err
}

func NewFileSystem(infile string) (Storage, error) {
	return &FileSystem{infile}, nil
}
