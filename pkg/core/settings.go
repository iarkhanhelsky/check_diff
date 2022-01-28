package core

import (
	"path/filepath"
	"strings"
)

type Settings struct {
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
		if match, _ := filepath.Match(pattern, r.File); match {
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
		if match, _ := filepath.Match(pattern, r.File); match {
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
