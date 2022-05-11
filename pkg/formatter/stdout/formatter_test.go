package stdout

import (
	"bytes"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestFormatter_Print(t *testing.T) {
	testCases := map[string]struct {
		issues   []core.Issue
		expected string
	}{
		"formatter_Main.java": {
			issues: []core.Issue{
				{
					File: "testdata/Main.java", Line: 6, Column: 9,
					Severity: "warn", Message: "Don't do that", Source: "jlint",
				},
			},
			expected: "testdata/formatter_Main.java.txt",
		},
		"no_issues": {
			issues:   []core.Issue{},
			expected: "testdata/no_issues.txt",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			expected, err := ioutil.ReadFile(tc.expected)
			assert.NoError(err)
			var buf bytes.Buffer
			NewFormatter().Print(tc.issues, &buf)
			assert.Equal(string(expected), buf.String())
		})
	}
}
