package checker

import (
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"path/filepath"
	"strings"
)

type fileChecker struct {
	supportedExtensions []string
	exclusions          []string
}

func (ch *fileChecker) filterExclusions(ranges []core.LineRange) []core.LineRange {
	size := 0
	for _, r := range ranges {
		if !matchExclusions(r, ch.exclusions) {
			ranges[size] = r
			size++
		}
	}
	return ranges[:size]
}

func matchExclusions(r core.LineRange, exclusions []string) bool {
	for _, e := range exclusions {

		if match, _ := filepath.Match(e, r.File); match {
			return true
		}
	}

	return false
}

func (ch *fileChecker) filterExtensions(ranges []core.LineRange) []core.LineRange {
	size := 0
	for _, r := range ranges {
		if matchExtensions(r, ch.supportedExtensions) {
			ranges[size] = r
			size++
		}
	}
	return ranges[:size]
}

func matchExtensions(r core.LineRange, extensions []string) bool {
	for _, e := range extensions {
		if strings.HasSuffix(r.File, e) {
			return true
		}
	}

	return false
}
