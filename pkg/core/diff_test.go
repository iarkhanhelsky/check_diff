package core

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func loadDiffLines(file string) ([]string, error) {
	diffFile, err := os.Open(path.Join("testdata", file))
	if err != nil {
		return nil, err
	}
	diffBytes, err := ioutil.ReadAll(diffFile)
	if err != nil {
		return nil, err
	}
	return strings.SplitN(string(diffBytes), "\n", -1), nil
}

func loadRanges(file string) ([]LineRange, error) {
	jsonFile, err := os.Open(path.Join("testdata", file))
	if err != nil {
		return nil, err
	}

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var ranges []LineRange
	if err := json.Unmarshal(jsonBytes, &ranges); err != nil {
		return nil, err
	}

	return ranges, nil
}

func TestDiffParser_ParseNextLine(t *testing.T) {
	assert := assert.New(t)
	for _, testCase := range []string{
		"single-file-one-line",
		"single-file-multiline",
		"single-file-multiple-range",
		"multiple-files-multiple-range",
	} {
		lines, err := loadDiffLines(testCase + ".diff")
		assert.NoError(err)
		expectRanges, err := loadRanges(testCase + ".json")
		assert.NoError(err)

		t.Run(testCase, func(t *testing.T) {
			parser := NewDiffParser()
			for _, line := range lines {
				parser.ParseNextLine(line)
			}
			assert.Equal(expectRanges, parser.Result())
		})
	}
}
