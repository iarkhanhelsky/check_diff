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

	issues, err := check.runChecks()
	if err != nil {
		return err
	}

	return check.writeIssues(issues)
}

func (check *Check) download() error {
	err := os.MkdirAll(check.Config.VendorDir, 0755)
	if err != nil {
		return err
	}
	for _, checker := range check.Checkers {
		for _, download := range checker.Downloads() {
			if err := download.Download(check.Config.VendorDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func (check *Check) runChecks() ([]core.Issue, error) {
	var issues []core.Issue

	issueChan := make(chan []core.Issue)
	errorChan := make(chan error)

	for _, checker := range check.Checkers {
		ch := checker
		go func() {
			r, err := ch.Check(check.LineRanges)
			if err != nil {
				errorChan <- err
			} else {
				issueChan <- r
			}
		}()
	}

	for sz := 0; sz < len(check.Checkers); sz++ {
		select {
		case i := <-issueChan:
			issues = append(issues, i...)
		case e := <-errorChan:
			return nil, fmt.Errorf("one or more checkers failed: %v", e)
		}
	}

	return issues, nil
}

func (check *Check) writeIssues(issues []core.Issue) error {
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
