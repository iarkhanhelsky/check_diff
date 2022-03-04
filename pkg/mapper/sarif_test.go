package mapper

import (
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/owenrumney/go-sarif/sarif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"path"
	"testing"
)

type SarifTestSuite struct {
	suite.Suite
}

func (s *SarifTestSuite) getBytes(file string) []byte {
	bytes, err := ioutil.ReadFile(file)
	assert.NoError(s.T(), err)
	return bytes
}

func (s *SarifTestSuite) TestParseSarif() {
	testcases := map[string]struct {
		report []byte
		issues []core.Issue
		err    bool
	}{
		"kube_linter_sarif_report.json": {
			report: s.getBytes("testdata/kube_linter_sarif_report.json"),
			issues: []core.Issue{
				{Tag: "no-read-only-root-fs", File: "k8s/deployment.yaml",
					Line: 1, Column: 0, Severity: "error",
					Message: "container \"sec-ctx-demo\" does not have a read-only root file system\nobject: <no namespace>/security-context-demo /v1, Kind=Pod",
					Source:  ""},
				{Tag: "unset-cpu-requirements", File: "k8s/deployment.yaml",
					Line: 1, Column: 0, Severity: "error",
					Message: "container \"sec-ctx-demo\" has cpu limit 0\nobject: <no namespace>/security-context-demo /v1, Kind=Pod",
					Source:  ""},
				core.Issue{Tag: "unset-memory-requirements", File: "k8s/deployment.yaml",
					Line: 1, Column: 0, Severity: "error",
					Message: "container \"sec-ctx-demo\" has memory limit 0\nobject: <no namespace>/security-context-demo /v1, Kind=Pod",
					Source:  ""}},
		},

		"bad report": {
			report: []byte("_"),
			err:    true,
		},
	}

	for name, tc := range testcases {
		s.T().Run(name, func(t *testing.T) {
			assert := assert.New(t)

			issues, err := parseSarif(tc.report)
			if tc.err {
				assert.Error(err)
			} else {
				assert.Equal(tc.issues, issues)
			}
		})
	}
}

func (s *SarifTestSuite) TestExtractLocation() {
	assert := s.Assert()

	testcases := []struct {
		report string
		run    int
		result int
		// expectations
		line   int
		column int
		file   string
	}{
		{
			report: "kube_linter_sarif_report.json", run: 0,
			line: 1, column: 0, file: "k8s/deployment.yaml",
		},
		{
			report: "checkstyle_sarif_report.json", run: 0, result: 3,
			line: 4, column: 5, file: "/Users/dm/Projects/github/check_diff/example/java/src/main/java/Main.java",
		},
	}

	for _, tc := range testcases {
		name := fmt.Sprintf("%s-%d-%d", tc.report, tc.run, tc.result)
		s.T().Run(name, func(t *testing.T) {
			report, _ := sarif.FromBytes(s.getBytes(path.Join("testdata", tc.report)))

			file, line, column := extractLocation(report.Runs[tc.run].Results[tc.result].Locations)

			assert.Equal(tc.file, file)
			assert.Equal(tc.line, line)
			assert.Equal(tc.column, column)
		})
	}
}

func TestReport(t *testing.T) {
	suite.Run(t, new(SarifTestSuite))
}
