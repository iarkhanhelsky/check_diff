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

	LineRanges []core.LineRange
	Checkers   []core.Checker `group:"checkers"`
	Formatter  core.Formatter
	Config     core.Config
}

type Check struct {
	LineRanges []core.LineRange
	Checkers   []core.Checker
	Formatter  core.Formatter
	Config     core.Config
}

func NewCheck(params Params) Check {
	return Check{
		LineRanges: params.LineRanges,
		Checkers:   params.Checkers,
		Formatter:  params.Formatter,
		Config:     params.Config,
	}
}

func (check *Check) Run() error {
	if err := check.download(); err != nil {
		return fmt.Errorf("failed to download dependencies: %v", err)
	}

	if err := check.runChecks(); err != nil {
		return fmt.Errorf("failed to run one or more checks: %v", err)
	}

	return nil
}

func (check *Check) download() error {
	for _, checker := range check.Checkers {
		for _, download := range checker.Downloads() {
			if err := download.Download(check.Config.VendorDir); err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func (check *Check) runChecks() error {
	var issues []core.Issue
	for _, checker := range check.Checkers {
		r, err := checker.Check(check.LineRanges)
		if err != nil {
			return fmt.Errorf("one or more checkers failed to finish: %v", err)
		}
		issues = append(issues, r...)
	}

	var writer io.Writer
	outFile := check.Config.OutputFile
	if len(outFile) == 0 {
		writer = os.Stdout
	} else {
		file, err := os.Create(outFile)
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
