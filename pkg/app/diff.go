package app

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func NewDiff(cliOptions CliOptions) ([]core.LineRange, error) {
	var reader io.Reader
	if len(cliOptions.InputFile) == 0 {
		reader = os.Stdin
	} else {
		file, err := os.Open(cliOptions.InputFile)
		if err != nil {
			return nil, fmt.Errorf("can't read file: %s: %v", cliOptions.InputFile, err)
		}
		defer file.Close()
		reader = file
	}
	diff, _ := ioutil.ReadAll(reader)
	parser := core.NewDiffParser()
	for _, line := range strings.Split(string(diff), "\n") {
		parser.ParseNextLine(line)
	}

	return parser.Result(), nil
}
