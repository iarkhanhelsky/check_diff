package core

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	assert := assert.New(t)

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

func TestMatch(t *testing.T) {
	type fixture struct {
		path    string
		pattern string
		matches bool
	}

	fixtures := []fixture{
		{"a/b/c/d", "a/b/**", false},
		{"a/b/c/d", "a/**/d", false},
		{"a/b/c/.d", "a/**/.d", false},
	}

	for _, f := range fixtures {
		t.Run(f.pattern, func(t *testing.T) {
			assert.Equal(t, f.matches, matchGlob(f.path, f.pattern))
		})
	}
}

func fromFiles(files ...string) []LineRange {
	var ranges []LineRange
	for _, f := range files {
		ranges = append(ranges, LineRange{File: f})
	}
	return ranges
}

func TestFilter(t *testing.T) {
	testcases := map[string]struct {
		settings   Settings
		ranges     []LineRange
		extensions []string
		expected   []LineRange
	}{
		"no filters": {
			settings: Settings{},
			ranges:   fromFiles("a.txt", "b.txt"),

			expected: fromFiles("a.txt", "b.txt"),
		},

		"filter extensions": {
			settings:   Settings{},
			ranges:     fromFiles("a.txt", "b.txt", "c.go"),
			extensions: []string{".go"},

			expected: fromFiles("c.go"),
		},

		"excludes": {
			settings: Settings{Exclude: []string{"gen/*"}},
			ranges:   fromFiles("gen/a", "gen/b.go", "c.go"),

			expected: fromFiles("c.go"),
		},

		"includes": {
			settings: Settings{Include: []string{"gen/*"}},
			ranges:   fromFiles("gen/a", "gen/b.go", "c.go"),

			expected: fromFiles("gen/a", "gen/b.go"),
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			assert.Equal(tc.expected, tc.settings.Filter(tc.ranges, tc.extensions...))
		})
	}
}

func TestSettings_CommandOrDefault(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("cat", Settings{Command: "cat"}.CommandOrDefault("xcat"))
	assert.Equal("xcat", Settings{Command: ""}.CommandOrDefault("xcat"))
}
