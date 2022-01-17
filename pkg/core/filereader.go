package core

import (
	"io/ioutil"
	"os"
)

type FileReader interface {
	Read(path string) (string, error)
}

type filereader struct {
}

var _ FileReader = &filereader{}

func (*filereader) Read(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	lines, err := ioutil.ReadAll(file)
	return (string(lines)), err
}