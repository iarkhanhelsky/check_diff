package app

import (
	"bufio"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"io"
	"os"
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

	parser := core.NewDiffParser()
	for scanner := bufio.NewScanner(reader); scanner.Scan(); {
		parser.ParseNextLine(scanner.Text())
	}

	return parser.Result(), nil
}
