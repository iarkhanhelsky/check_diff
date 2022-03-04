package stdout

import (
	"bytes"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestFormatter_Print(t *testing.T) {
	assert := assert.New(t)
	issues := []core.Issue{
		{
			File: "testdata/Main.java", Line: 6, Column: 9,
			Severity: "warn", Message: "Don't do that", Source: "jlint",
		},
	}

	expectedOutput, err := ioutil.ReadFile("testdata/formatter_Main.java.txt")
	assert.NoError(err)

	var buf bytes.Buffer
	NewFormatter().Print(issues, &buf)
	assert.Equal(string(expectedOutput), buf.String())
}
