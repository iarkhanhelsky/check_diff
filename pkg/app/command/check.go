package command

import (
	"bufio"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/tools"
	"io"
	"os"
	"time"

	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"go.uber.org/zap"
)

type Check struct {
	Env      Env
	Checkers []core.Checker
	Config   core.Config
	Registry tools.Registry
	Logger   *zap.SugaredLogger
}

var _ Command = &Check{}

func NewCheck(env Env, config core.Config, checkers []core.Checker, logger *zap.SugaredLogger, registry tools.Registry) Command {
	return &Check{
		Env:      env,
		Config:   config,
		Checkers: checkers,
		Registry: registry,
		Logger:   logger,
	}
}

func (check *Check) Run() error {
	var err error
	if err = check.download(); err != nil {
		return fmt.Errorf("downloading dependencies: %v", err)
	}

	var ranges []core.LineRange
	if ranges, err = check.readDiff(); err != nil {
		return fmt.Errorf("reading diff: %v", err)
	}

	issues, err := check.runChecks(ranges)
	if err != nil {
		return err
	}

	return check.writeIssues(issues)
}

func (check *Check) readDiff() ([]core.LineRange, error) {
	var reader io.Reader

	config := check.Config
	logger := check.Logger
	if len(config.InputFile) == 0 {
		logger.Debugf("reading diff from STDIN")
		reader = os.Stdin
	} else {
		file, err := os.Open(config.InputFile)
		logger.With("file", config.InputFile).Debugf("reading diff from file")
		if err != nil {
			return nil, fmt.Errorf("can't read file: %s: %v", config.InputFile, err)
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

func (check *Check) download() error {
	err := os.MkdirAll(check.Config.VendorDir, 0755)
	if err != nil {
		return fmt.Errorf("preparing vendor dir %s: %w", check.Config.VendorDir, err)
	}

	var binaries []*tools.Binary
	for _, checker := range check.Checkers {
		if c, ok := checker.(core.Binaries); ok {
			binaries = append(binaries, c.Binaries()...)
		}
	}

	return check.Registry.Install(binaries...)
}

func (check *Check) runChecks(lineRanges []core.LineRange) ([]core.Issue, error) {
	var issues []core.Issue

	issueChan := make(chan []core.Issue)
	errorChan := make(chan error)

	for _, checker := range check.Checkers {
		ch := checker
		go func() {
			start := time.Now()
			r, err := ch.Check(lineRanges)
			duration := time.Since(start)

			if err != nil {
				check.Logger.With("checker", ch.Tag()).Errorf("finished with error: %v", err)
				errorChan <- err
			} else {
				issueChan <- r
			}
			check.Logger.With("checker", ch.Tag()).Debugf("took %s", duration)
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
		writer = check.Env.OutWriter
	} else {
		file, err := os.Create(outFile)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("can't open output file: %s: %v", outFile, err)
		}
		writer = file
	}

	formatter, err := formatter.NewFormatter(formatter.Params{Format: check.Config.OutputFormat})
	if err != nil {
		return fmt.Errorf("can't create formatter: %v", err)
	}
	if err := formatter.Print(issues, writer); err != nil {
		return fmt.Errorf("can't print issues: %v", err)
	}

	return nil
}
