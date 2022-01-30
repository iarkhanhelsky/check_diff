package command

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/fx"
	"io"
	"os"
)

type Params struct {
	fx.In
}

type Check struct {
	LineRanges []core.LineRange
	Checker    core.Checker
	Formatter  core.Formatter
	Config     core.Config
}

func NewCheck(lineRanges []core.LineRange,
	checker core.Checker,
	formatter core.Formatter,
	config core.Config) Check {
	return Check{
		LineRanges: lineRanges,
		Checker:    checker,
		Formatter:  formatter,
		Config:     config,
	}
}

func (check *Check) Run() error {
	issues, err := check.Checker.Check(check.LineRanges)
	if err != nil {
		return fmt.Errorf("one or more checkers failed to finish: %v", err)
	}

	var writer io.Writer
	outFile := check.Config.OutputFile
	if len(outFile) == 0 {
		writer = os.Stdout
	} else {
		file, err := os.Open(outFile)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("can't open output file: %s: %v", outFile, err)
		}
		writer = file
	}

	if err := check.Formatter.Print(issues, writer); err != nil {
		return fmt.Errorf("can't print issues: %v", err)
	}

	return nil
}
