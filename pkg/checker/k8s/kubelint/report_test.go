package kubelint

import (
	"github.com/owenrumney/go-sarif/sarif"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type ReportTestSuite struct {
	suite.Suite

	reportBytes []byte
}

func (s *ReportTestSuite) SetupTest() {
	file, err := os.Open("testdata/sarif_report.json")
	defer file.Close()
	assert.NoError(s.T(), err)
	bytes, err := ioutil.ReadAll(file)
	assert.NoError(s.T(), err)

	s.reportBytes = bytes
}

func (s *ReportTestSuite) TestExtractLocation() {
	assert := assert.New(s.T())

	report, _ := sarif.FromBytes(s.reportBytes)

	file, line, column := extractLocation(report.Runs[0].Results[0].Locations)
	assert.Equal("k8s/deployment.yaml", file)
	assert.Equal(1, line)
	assert.Equal(0, column)
}

func TestReport(t *testing.T) {
	suite.Run(t, new(ReportTestSuite))
}
