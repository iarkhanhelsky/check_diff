package integration_tests

import (
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestGitlab(t *testing.T) {
	testCommon(t, testCases(t), formatter.Gitlab)
}

func TestPhabricator(t *testing.T) {
	testCommon(t, testCases(t), formatter.Phabricator)
}

func TestStdout(t *testing.T) {
	testCommon(t, testCases(t), formatter.STDOUT)
}

func testCommon(t *testing.T, inputs []string, format formatter.Format) {

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
