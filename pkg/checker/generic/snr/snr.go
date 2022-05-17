package snr

import (
	"errors"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"go.uber.org/config"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ScriptAndRegexp ...
//
// The regexp should capture these named patterns with (?P<name>...):
//  * `tag`      string
//  * `file`     string
//  * `line`     int
//  * `column`   int
//  * `severity` string
//  * `message`  string
//  * `source`   string
type ScriptAndRegexp struct {
	Enabled bool     `yaml:"Enabled"`
	Exclude []string `yaml:"Exclude"`
	Include []string `yaml:"Include"`
	Script  string   `yaml:"Script"`
	Regexp  string   `yaml:"Regexp"`
	regexp  *regexp.Regexp
	tag     string
}

var _ core.Checker = &ScriptAndRegexp{}

func (snr ScriptAndRegexp) Tag() string {
	return snr.tag
}

func (snr ScriptAndRegexp) Check(ranges []core.LineRange) ([]core.Issue, error) {
	flow := core.NewFlow(snr.Tag(), core.Settings{
		Enabled: true,
		Exclude: snr.Exclude,
		Include: snr.Include,
		Command: snr.Script,
	}, core.WithConverter(snr.regexpConverter),
		core.WithCommand(snr.Script))

	return checkAll(ranges, flow)
}

func checkAll(ranges []core.LineRange, flow core.Flow) ([]core.Issue, error) {
	rangesCopy := make([]core.LineRange, len(ranges))
	copy(rangesCopy, ranges)
	// Run once per file
	sort.Slice(rangesCopy, func(i, j int) bool {
		return strings.Compare(rangesCopy[i].File, rangesCopy[j].File) < 0
	})

	begin := 0
	result := []core.Issue{}
	for i, rng := range rangesCopy {
		if begin != i && rangesCopy[begin].File != rng.File {
			out, err := flow.Run(rangesCopy[begin:i])
			if err != nil {
				return []core.Issue{}, fmt.Errorf("checking %s: %w", rangesCopy[begin].File, err)
			}
			result = append(result, out...)
			begin = i
		}
	}
	if begin < len(rangesCopy) {
		out, err := flow.Run(rangesCopy[begin:])
		if err != nil {
			return []core.Issue{}, fmt.Errorf("checking %s: %w", rangesCopy[begin].File, err)
		}
		result = append(result, out...)
	}
	return result, nil
}

func (snr ScriptAndRegexp) regexpConverter(bytes []byte) ([]core.Issue, error) {
	lines := strings.Split(string(bytes), "\n")
	out := []core.Issue{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		issue := core.Issue{}
		match := snr.regexp.FindStringSubmatch(line)
		for i, name := range snr.regexp.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}

			if err := snr.updateIssueField(&issue, name, match[i]); err != nil {
				return nil, fmt.Errorf("parsing %snr output: %w", snr.tag, err)
			}
		}

		out = append(out, issue)
	}
	return out, nil
}

func (snr ScriptAndRegexp) updateIssueField(issue *core.Issue, name string, value string) error {
	if name == "file" {
		issue.File = value
	} else if name == "tag" {
		issue.Tag = value
	} else if name == "line" {
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("line=%s: %w", value, err)
		}
		issue.Line = v
	} else if name == "column" {
		v, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("column=%s: %w", value, err)
		}
		issue.Column = v
	} else if name == "severity" {
		issue.Severity = value
	} else if name == "message" {
		issue.Message = value
	} else if name == "source" {
		issue.Source = value
	}

	return nil
}

func ScriptAndRegexProvider(yaml *config.YAML) ([]core.Checker, error) {
	raw := make(map[string]interface{})
	err := yaml.Get("").Populate(&raw)
	if err != nil {
		return nil, fmt.Errorf("reading config for ScriptAndRegexp linters: %w", err)
	}
	var checkers []core.Checker
	for k, _ := range raw {
		if strings.HasPrefix(k, "ScriptAndRegexp/") {
			snr := &ScriptAndRegexp{tag: k}
			if err := yaml.Get(k).Populate(snr); err != nil {
				return nil, fmt.Errorf("reading config for %s: %w", k, err)
			}

			if snr.Regexp == "" {
				return nil, fmt.Errorf("reading config for %s: Regexp is empty", k)
			}
			rx, err := regexp.Compile(snr.Regexp)
			if err != nil {
				return nil, fmt.Errorf("reading config for %s: %s is invalid regexp: %w", k, snr.Regexp, err)
			}
			snr.regexp = rx

			if snr.Script == "" {
				return nil, fmt.Errorf("reading config for %s: Script is empty", k)
			}
			if _, err := os.Stat(snr.Script); errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("reading config for %s: Script %s does not exist", k, snr.Script)
			}

			checkers = append(checkers, snr)
		}
	}

	return checkers, nil
}
