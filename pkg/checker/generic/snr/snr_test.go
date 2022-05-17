package snr

import (
	"github.com/golang/mock/gomock"
	mockcore "github.com/iarkhanhelsky/check_diff/mocks/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/stretchr/testify/assert"
	"go.uber.org/config"
	"os"
	"regexp"
	"testing"
)

func TestScriptAndRegexProvider(t *testing.T) {
	testCases := map[string]struct {
		checkers []core.Checker
		fixture  string
		err      string
	}{
		"invalid_param": {
			fixture: "testdata/snr/invalid_param.yaml",
			err:     "reading config for ScriptAndRegexp/test: yaml: unmarshal errors:\n  line 2: cannot unmarshal !!bool `false` into []string",
		},
		"single": {
			fixture: "testdata/snr/single.yaml",
			checkers: []core.Checker{
				&ScriptAndRegexp{
					tag:     "ScriptAndRegexp/test",
					Enabled: true,
					Script:  "testdata/snr/lint.sh",
					Regexp:  "(?P<file>[^:]*):(?P<line>[0-9]*):(?P<message>.*)$$",
					regexp:  regexp.MustCompile("(?P<file>[^:]*):(?P<line>[0-9]*):(?P<message>.*)$$"),
				},
			},
		},
		"multiple": {
			fixture: "testdata/snr/multiple.yaml",
			checkers: []core.Checker{
				&ScriptAndRegexp{
					tag:     "ScriptAndRegexp/test1",
					Enabled: true,
					Script:  "testdata/snr/lint.sh",
					Regexp:  "(?P<file>[^:]*):(?P<line>[0-9]*):(?P<message>.*)$$",
					regexp:  regexp.MustCompile("(?P<file>[^:]*):(?P<line>[0-9]*):(?P<message>.*)$$"),
				},
				&ScriptAndRegexp{
					tag:     "ScriptAndRegexp/test2",
					Enabled: true,
					Script:  "testdata/snr/lint.sh",
					Regexp:  "(?P<file>[^:]*):(?P<line>[0-9]*):(?P<message>.*)$$",
					regexp:  regexp.MustCompile("(?P<file>[^:]*):(?P<line>[0-9]*):(?P<message>.*)$$"),
				},
			},
		},
		"invalid regexp": {
			fixture: "testdata/snr/invalid_regexp.yaml",
			err:     "reading config for ScriptAndRegexp/test: [ is invalid regexp: error parsing regexp: missing closing ]: `[`",
		},
		"script does not exist": {
			fixture: "testdata/snr/does_not_exist.yaml",
			err:     "reading config for ScriptAndRegexp/test: Script testdata/snr/lint1.sh does not exist",
		},
		"script is empty": {
			fixture: "testdata/snr/script_is_empty.yaml",
			err:     "reading config for ScriptAndRegexp/test: Script is empty",
		},
		"regexp is empty": {
			fixture: "testdata/snr/regexp_is_empty.yaml",
			err:     "reading config for ScriptAndRegexp/test: Regexp is empty",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			file, err := os.Open(tc.fixture)
			defer file.Close()
			assert.NoError(err)

			yaml, err := config.NewYAML(config.RawSource(file))
			assert.NoError(err)

			checkers, err := ScriptAndRegexProvider(yaml)
			if tc.err != "" {
				assert.EqualError(err, tc.err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.checkers, checkers)
			}
		})
	}
}

func TestScriptAndRegexp_CheckAll(t *testing.T) {
	testCases := map[string]struct {
		input []core.LineRange
		calls [][]core.LineRange
	}{
		"single": {
			input: []core.LineRange{
				{File: "test1.go", Start: 5, End: 6},
			},
			calls: [][]core.LineRange{
				{
					{File: "test1.go", Start: 5, End: 6},
				},
			},
		},
		"double": {
			input: []core.LineRange{
				{File: "test1.go", Start: 5, End: 6},
				{File: "test1.go", Start: 6, End: 7},
			},
			calls: [][]core.LineRange{
				{
					{File: "test1.go", Start: 5, End: 6},
					{File: "test1.go", Start: 6, End: 7},
				},
			},
		},
		"two groups of single": {
			input: []core.LineRange{
				{File: "test1.go", Start: 5, End: 6},
				{File: "test2.go", Start: 6, End: 7},
			},
			calls: [][]core.LineRange{
				{
					{File: "test1.go", Start: 5, End: 6},
				},
				{
					{File: "test2.go", Start: 6, End: 7},
				},
			},
		},
		"two groups of double": {
			input: []core.LineRange{
				{File: "test1.go", Start: 5, End: 6},
				{File: "test1.go", Start: 6, End: 7},
				{File: "test2.go", Start: 5, End: 6},
				{File: "test2.go", Start: 6, End: 7},
			},
			calls: [][]core.LineRange{
				{
					{File: "test1.go", Start: 5, End: 6},
					{File: "test1.go", Start: 6, End: 7},
				},
				{
					{File: "test2.go", Start: 5, End: 6},
					{File: "test2.go", Start: 6, End: 7},
				},
			},
		},
		"three groups of double": {
			input: []core.LineRange{
				{File: "test1.go", Start: 5, End: 6},
				{File: "test1.go", Start: 6, End: 7},
				{File: "test2.go", Start: 5, End: 6},
				{File: "test2.go", Start: 6, End: 7},
				{File: "test3.go", Start: 5, End: 6},
				{File: "test3.go", Start: 6, End: 7},
			},
			calls: [][]core.LineRange{
				{
					{File: "test1.go", Start: 5, End: 6},
					{File: "test1.go", Start: 6, End: 7},
				},
				{
					{File: "test2.go", Start: 5, End: 6},
					{File: "test2.go", Start: 6, End: 7},
				},
				{
					{File: "test3.go", Start: 5, End: 6},
					{File: "test3.go", Start: 6, End: 7},
				},
			},
		},
		"unsorted": {
			input: []core.LineRange{
				{File: "test5.go", Start: 5, End: 6},
				{File: "test8.go", Start: 6, End: 7},
				{File: "test1.go", Start: 5, End: 6},
				{File: "test2.go", Start: 6, End: 7},
				{File: "test4.go", Start: 5, End: 6},
				{File: "test6.go", Start: 6, End: 7},
			},
			calls: [][]core.LineRange{
				{{File: "test1.go", Start: 5, End: 6}},
				{{File: "test2.go", Start: 6, End: 7}},
				{{File: "test4.go", Start: 5, End: 6}},
				{{File: "test5.go", Start: 5, End: 6}},
				{{File: "test6.go", Start: 6, End: 7}},
				{{File: "test8.go", Start: 6, End: 7}},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			ctrl := gomock.NewController(t)
			flow := mockcore.NewMockFlow(ctrl)
			for _, args := range tc.calls {
				flow.EXPECT().Run(args).Return([]core.Issue{}, nil)
			}

			out, err := checkAll(tc.input, flow)
			assert.NoError(err)
			assert.Empty(out)
		})
	}
}

func TestScriptAndRegexp_regexpConverter(t *testing.T) {
	testCases := map[string]struct {
		input    string
		regexp   *regexp.Regexp
		expected []core.Issue
		err      string
	}{
		"empty": {
			input:    "",
			expected: []core.Issue{},
		},
		"file line column message": {
			input:  "test.go:15:4:Undeclared variable boo",
			regexp: regexp.MustCompile(`(?P<file>[^:]*):(?P<line>[0-9]*):(?P<column>[0-9]*):(?P<message>.*)$`),
			expected: []core.Issue{
				{
					File:    "test.go",
					Line:    15,
					Column:  4,
					Message: "Undeclared variable boo",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			snr := ScriptAndRegexp{
				regexp: tc.regexp,
			}
			out, err := snr.regexpConverter([]byte(tc.input))
			if tc.err != "" {
				assert.EqualError(err, tc.err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.expected, out)
			}
		})
	}
}

func TestScriptAndRegexp_UpdateIssueField(t *testing.T) {
	testCases := map[string]struct {
		name   string
		value  string
		expect core.Issue
		err    string
	}{
		"message": {
			name:  "message",
			value: "foo bar",
			expect: core.Issue{
				Message: "foo bar",
			},
		},
		"file": {
			name:  "file",
			value: "test.go",
			expect: core.Issue{
				File: "test.go",
			},
		},
		"tag": {
			name:  "tag",
			value: "cool-lint",
			expect: core.Issue{
				Tag: "cool-lint",
			},
		},
		"line": {
			name:  "line",
			value: "10",
			expect: core.Issue{
				Line: 10,
			},
		},
		"line invalid": {
			name:  "line",
			value: "no",
			err:   "line=no: strconv.Atoi: parsing \"no\": invalid syntax",
		},
		"column": {
			name:  "column",
			value: "10",
			expect: core.Issue{
				Column: 10,
			},
		},
		"column invalid": {
			name:  "column",
			value: "no",
			err:   "column=no: strconv.Atoi: parsing \"no\": invalid syntax",
		},
		"severity": {
			name:  "severity",
			value: "warning",
			expect: core.Issue{
				Severity: "warning",
			},
		},
		"source": {
			name:  "source",
			value: "LintRule100",
			expect: core.Issue{
				Source: "LintRule100",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			issue := core.Issue{}
			err := ScriptAndRegexp{}.updateIssueField(&issue, tc.name, tc.value)
			if tc.err != "" {
				assert.EqualError(err, tc.err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.expect, issue)
			}
		})
	}
}
