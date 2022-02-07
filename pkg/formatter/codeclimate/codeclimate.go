package codeclimate

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"io"
	"strconv"
)

// See https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md
// Example:
// ```
//     [
//       {
//         "description": "'unused' is assigned a value but never used.",
//         "fingerprint": "7815696ecbf1c96e6894b779456d330e",
//         "severity": "minor",
//         "location": {
//           "path": "lib/index.js",
//           "lines": {
//             "begin": 42
//           }
//         }
//       }
//     ]
// ```
type Formatter struct{}

type issue struct {
	Description string   `json:"description"`
	Fingerprint string   `json:"fingerprint"`
	Severity    string   `json:"severity"`
	Location    location `json:"location"`
}

type location struct {
	Path  string `json:"path"`
	Lines lines  `json:"lines"`
}

type lines struct {
	Begin int `json:"begin"`
}

var _ core.Formatter = &Formatter{}

func (cc *Formatter) Supports() []core.Format {
	return []core.Format{core.Codeclimate, core.Gitlab}
}

func (cc *Formatter) Print(issues []core.Issue, w io.Writer) error {
	var report []issue
	for _, issue := range issues {
		e := convert(issue)
		report = append(report, e)
	}

	bytes, err := json.Marshal(&report)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

func convert(iss core.Issue) issue {
	return issue{
		Description: iss.Message,
		Fingerprint: computeFingerprint(iss),
		Severity:    mapSeverity(iss.Severity),
		Location: location{
			Path: iss.File,
			Lines: lines{
				Begin: iss.Line,
			},
		},
	}
}

func computeFingerprint(issue core.Issue) string {
	var buffer bytes.Buffer
	buffer.WriteString(issue.File)
	buffer.WriteString(strconv.Itoa(issue.Line))
	buffer.WriteString(strconv.Itoa(issue.Column))
	buffer.WriteString(issue.Severity)

	return fmt.Sprintf("%x", md5.Sum(buffer.Bytes()))
}

func mapSeverity(severity string) string {
	return severity
}
