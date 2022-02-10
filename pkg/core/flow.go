package core

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Flow interface {
	Run(ranges []LineRange) ([]Issue, error)
}

type FlowOption func(*flow)

func WithCommand(command string, args ...string) FlowOption {
	return func(f *flow) {
		f.args = args
		f.command = command
	}
}

func WithConverter(converter Converter) FlowOption {
	return func(f *flow) {
		f.converter = converter
	}
}

func WithFileExtensions(exts ...string) FlowOption {
	return func(f *flow) {
		f.supportedExtensions = exts
	}
}

type Converter func([]byte) ([]Issue, error)

func NewFlow(tag string, s Settings, opts ...FlowOption) Flow {
	f := flow{tag: tag, settings: s}
	for _, o := range opts {
		o(&f)
	}

	return &f
}

type flow struct {
	tag                 string
	settings            Settings
	command             string
	args                []string
	converter           Converter
	supportedExtensions []string
}

var _ Flow = &flow{}

func (f *flow) Run(ranges []LineRange) ([]Issue, error) {
	matchedRanges := f.settings.Filter(ranges, f.supportedExtensions...)
	if len(matchedRanges) == 0 {
		return []Issue{}, nil
	}

	args := []string{}
	args = append(args, f.args...)

	for _, r := range matchedRanges {
		args = append(args, r.File)
	}

	cmd := exec.Command(f.command, args...)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil && cmd.ProcessState.ExitCode() != 1 {
		return nil, fmt.Errorf("failed to run %s: %v: %s", f.tag, err, string(stderr.Bytes()))
	}

	issues, err := f.converter(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s issues: %v", f.tag, err)
	}

	sz := 0
	for _, issue := range issues {
		matched := false
		for _, r := range ranges {
			if r.File == issue.File && r.Start <= issue.Line && issue.Line <= r.End {
				matched = true
				break
			}
		}
		if matched {
			issues[sz] = issue
			sz++
		}
	}
	return issues[:sz], nil
}
