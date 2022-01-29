package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"testing"
)

type SettingsTestSuite struct {
	suite.Suite
}

func (s *SettingsTestSuite) TestUnmarshal() {
	assert := assert.New(s.T())

	raw := `
Exclude:
- gen/*.java
Include:
- src/main/java/**
Command: check_java
Config: check_java.xml
`
	opts := Settings{}
	err := yaml.Unmarshal([]byte(raw), &opts)
	assert.NoError(err)
	assert.Equal(Settings{
		Exclude: []string{"gen/*.java"},
		Include: []string{"src/main/java/**"},
		Command: "check_java",
		Config:  "check_java.xml",
	}, opts)
}

func (s *SettingsTestSuite) TestMatch() {
	type fixture struct {
		path    string
		pattern string
		matches bool
	}

	fixtures := []fixture{
		{"a/b/c/d", "a/b/**", true},
		{"a/b/c/d", "a/**/d", true},
		{"a/b/c/.d", "a/**/.d", true},
	}

	for _, f := range fixtures {
		s.T().Run(f.pattern, func(t *testing.T) {
			assert.Equal(t, f.matches, matchGlob(f.path, f.pattern))
		})
	}
}

func TestOptions(t *testing.T) {
	suite.Run(t, new(SettingsTestSuite))
}
