package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"testing"
)

type OptionsTestSuite struct {
	suite.Suite
}

func (s *OptionsTestSuite) TestUnmarshal() {
	assert := assert.New(s.T())

	raw := `
Exclude:
- gen/*.java
Include:
- src/main/java/**
Command: check_java
Config: check_java.xml
`
	opts := Options{}
	err := yaml.Unmarshal([]byte(raw), &opts)
	assert.NoError(err)
	assert.Equal(Options{
		Exclude: []string{"gen/*.java"},
		Include: []string{"src/main/java/**"},
		Command: "check_java",
		Config:  "check_java.xml",
	}, opts)
}

func TestOptions(t *testing.T) {
	suite.Run(t, new(OptionsTestSuite))
}
