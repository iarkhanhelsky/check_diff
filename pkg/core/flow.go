package core

import (
	"fmt"
)

type Flow interface {
	Run(ranges []LineRange) ([]Issue, error)
}

type FlowOption func(*flow)
type Converter func([]byte) ([]Issue, error)
type ArgFunction func(lineRange LineRange) []string

var DefaultArgFunction ArgFunction = func(lineRange LineRange) []string {
	return []string{lineRange.File}
}

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

func WithShellOptions(opts ...ShellOption) FlowOption {
	return func(f *flow) {
		f.shellOptions = opts
	}
}

func WithShellFactory(factory func(...ShellOption) Shell) FlowOption {
	return func(f *flow) {
		f.shellFactory = factory
	}
}

func NewFlow(tag string, s Settings, opts ...FlowOption) Flow {
	f := flow{
		tag:          tag,
		settings:     s,
		argFunction:  DefaultArgFunction,
		shellOptions: []ShellOption{AllowExitCodes(0, 1)},
		shellFactory: NewLocalShell,
	}
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
	argFunction         ArgFunction
	shellOptions        []ShellOption
	converter           Converter
	supportedExtensions []string
	shellFactory        func(opts ...ShellOption) Shell
}

var _ Flow = &flow{}

func (f *flow) Run(ranges []LineRange) ([]Issue, error) {
	matchedRanges := f.settings.Filter(ranges, f.supportedExtensions...)
	if len(matchedRanges) == 0 {
		return []Issue{}, nil
	}

	args := f.args
	for _, r := range matchedRanges {
		args = append(args, f.argFunction(r)...)
	}

	output, err := f.shellFactory(f.shellOptions...).Execute(f.command, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s issues: %v", f.tag, err)
	}

	issues, err := f.converter(output)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s issues: %v", f.tag, err)
	}

	issues = filterIssues(ranges, issues)

	return issues, nil
}

func filterIssues(ranges []LineRange, issues []Issue) []Issue {
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

	return issues[:sz]
}
