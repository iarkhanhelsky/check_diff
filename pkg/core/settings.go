package core

import (
	fnmatch "github.com/iarkhanhelsky/fnmatch.v2"
	"strings"
)

type Settings struct {
	Enabled bool     `yaml:"Enabled"`
	Exclude []string `yaml:"Exclude"`
	Include []string `yaml:"Include"`
	Command string   `yaml:"Command"`
	Config  string   `yaml:"Config"`
}

func (settings Settings) Filter(ranges []LineRange, supportedExtensions ...string) []LineRange {
	var result []LineRange
	for _, r := range ranges {
		if matchesExtensions(r, supportedExtensions) && settings.isIncluded(r) && !settings.isExcluded(r) {
			result = append(result, r)
		}
	}
	return result
}

func (settings Settings) isIncluded(r LineRange) bool {
	if len(settings.Include) == 0 {
		return true
	}

	for _, pattern := range settings.Include {
		if matchGlob(pattern, r.File) {
			return true
		}
	}

	return false
}

func (settings Settings) isExcluded(r LineRange) bool {
	if len(settings.Include) == 0 {
		return false
	}

	for _, pattern := range settings.Exclude {
		if matchGlob(pattern, r.File) {
			return true
		}
	}

	return false
}

func matchesExtensions(r LineRange, extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}

	for _, e := range extensions {
		if strings.HasSuffix(r.File, e) {
			return true
		}
	}

	return false
}

func matchGlob(pattern string, path string) bool {
	return fnmatch.Match(pattern, path, 0)
}
