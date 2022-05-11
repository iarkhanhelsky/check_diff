package golangci_lint

import (
	"encoding/json"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/tools"
	"strings"
)

type Checker struct {
	core.Settings `yaml:",inline"`

	golangciLint *tools.Binary
}

var _ core.Checker = &Checker{}
var _ core.Binaries = &Checker{}

const (
	golangciLintVersion     = "1.46.0"
	golangciLintUrlTemplate = "https://github.com/golangci/golangci-lint/releases/download/v@VERSION/golangci-lint-@VERSION-@TARGET.tar.gz"
)

func NewChecker() core.Checker {
	return &Checker{
		golangciLint: &tools.Binary{
			Name:    "golangci-lint",
			DstFile: "golangci-lint.tar.gz",
			Path:    strings.ReplaceAll("golangci-lint-@VERSION-darwin-arm64/golangci-lint", "@VERSION", golangciLintVersion),
			Targets: map[tools.TargetTuple]tools.TargetSource{
				tools.DarwinAMD64: {
					Urls: []string{strings.ReplaceAll(
						strings.ReplaceAll(golangciLintUrlTemplate, "@VERSION", golangciLintVersion), "@TARGET", "darwin-amd64")},
					SHA256: "a9a502ba8789746e8f667398ba64cc3b7d5a2e8b112d6704ba13c47ec7e8ec3b",
				},
				tools.DarwinARM64: {
					Urls: []string{strings.ReplaceAll(
						strings.ReplaceAll(golangciLintUrlTemplate, "@VERSION", golangciLintVersion), "@TARGET", "darwin-arm64")},
					SHA256: "7c6aab18341444ab47281e68088582a353102a3e7b637be8d6466190d720d100",
				},
				tools.LinuxAMD64: {
					Urls: []string{strings.ReplaceAll(
						strings.ReplaceAll(golangciLintUrlTemplate, "@VERSION", golangciLintVersion), "@TARGET", "linux-amd64")},
					SHA256: "7f50c49674660baaddd8b6ce1c815563bb1dd238805518b11d3021e613481927",
				},
				tools.LinuxARM64: {
					Urls: []string{strings.ReplaceAll(
						strings.ReplaceAll(golangciLintUrlTemplate, "@VERSION", golangciLintVersion), "@TARGET", "linux-arm64")},
					SHA256: "129bf5665727c0844b047339c4b54091f87bd07227311a4b2cd255d51c9748dd",
				},
			},
		},
	}
}

func (checker *Checker) Tag() string {
	return "golangci-lint"
}

func (checker *Checker) Check(ranges []core.LineRange) ([]core.Issue, error) {
	args := []string{"run", "--out-format", "json"}
	if checker.Settings.Config != "" {
		args = append(args, "-c", checker.Settings.Config)
	}
	return core.NewFlow(checker.Tag(), checker.Settings,
		core.WithCommand(checker.golangciLint.Executable(), args...),
		core.WithFileExtensions(".go"),
		core.WithConverter(convertLintJson),
	).Run(ranges)
}

func (checker *Checker) Binaries() []*tools.Binary {
	return []*tools.Binary{checker.golangciLint}
}

func convertLintJson(bytes []byte) ([]core.Issue, error) {
	var report report
	if err := json.Unmarshal(bytes, &report); err != nil {
		return nil, fmt.Errorf("converting issues: %w", err)
	}
	var result []core.Issue
	for _, issue := range report.Issues {
		severity := issue.Severity
		if issue.Severity == "" {
			severity = "error"
		}
		result = append(result, core.Issue{
			File:     issue.Pos.Filename,
			Line:     issue.Pos.Line,
			Column:   issue.Pos.Column,
			Severity: severity,
			Message:  issue.Text,
			Source:   issue.FromLinter,
			Tag:      issue.FromLinter,
		})
	}
	return result, nil
}

type report struct {
	Issues []issue `json:"Issues"`
}

type issue struct {
	FromLinter string `json:"FromLinter"`
	Text       string `json:"Text"`
	Severity   string `json:"Severity"`
	Pos        pos    `json:"Pos"`
}

type pos struct {
	Filename string `json:"Filename"`
	Line     int    `json:"Line"`
	Column   int    `json:"Column"`
}
