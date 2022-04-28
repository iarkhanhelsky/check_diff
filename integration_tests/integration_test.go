package integration_tests

import (
	"bytes"
	"github.com/iarkhanhelsky/check_diff/pkg/app"
	"github.com/iarkhanhelsky/check_diff/pkg/app/command"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestGitlab(t *testing.T) {
	testCommon(t, testCases(t), formatter.Gitlab, ".gitlab.json")
}

func TestPhabricator(t *testing.T) {
	testCommon(t, testCases(t), formatter.Phabricator, ".phabricator.json")
}

func TestStdout(t *testing.T) {
	testCommon(t, testCases(t), formatter.STDOUT, ".stdout.txt")
}

func testCommon(t *testing.T, inputs []string, format formatter.Format, suffix string) {
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			data, err := ioutil.ReadFile(input + suffix)
			assert.NoError(t, err)
			testCmd(t, data, "check_diff", "-i", input, "-f", string(format))
		})
	}
}

func testCmd(t *testing.T, output []byte, args ...string) {
	var outBuff bytes.Buffer
	var errBuff bytes.Buffer

	err := app.Main(command.Env{
		Args:      append(args, "--vendor-dir", t.TempDir(), "--no-color"),
		OutWriter: &outBuff,
		ErrWriter: &errBuff,
	})
	assert.NoError(t, err)
	assert.Equal(t, "", errBuff.String())
	assert.Equal(t, string(output), outBuff.String())
}

func testCases(t *testing.T) []string {
	entries, err := ioutil.ReadDir("diffs")
	assert.NoError(t, err)

	var cases []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".diff") {
			cases = append(cases, path.Join("diffs", e.Name()))
		}
	}

	return cases
}
