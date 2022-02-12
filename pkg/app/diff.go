package app

import (
	"bufio"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"os"
)

type DiffParams struct {
	fx.In

	CliOptions CliOptions
	Logger     *zap.Logger
}

func NewDiff(params DiffParams) ([]core.LineRange, error) {
	var reader io.Reader

	cliOptions := params.CliOptions
	logger := params.Logger.Sugar()
	if len(cliOptions.InputFile) == 0 {
		logger.Debugf("reading diff from STDIN")
		reader = os.Stdin
	} else {
		file, err := os.Open(cliOptions.InputFile)
		logger.With("file", cliOptions.InputFile).Debugf("reading diff from file")
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
