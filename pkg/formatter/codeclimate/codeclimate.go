package codeclimate

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/iarkhanhelsky/check_diff/pkg/core"
	"github.com/iarkhanhelsky/check_diff/pkg/formatter"
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
	description string   `json:"description"`
	fingerprint string   `json:"fingerprint"`
	severity    string   `json:"severity"`
	location    location `json:"location"`
}

type location struct {
	path  string `json:"path"`
	lines lines  `json:"lines"`
}

type lines struct {
	begin int `json:"begin"`
}

var _ formatter.Formatter = &Formatter{}

func (cc *Formatter) Supports() []formatter.Format {
	return []formatter.Format{formatter.Codeclimate, formatter.Gitlab}
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
		description: iss.Message,
		fingerprint: computeFingerprint(iss),
		severity:    mapSeverity(iss.Severity),
		location: location{
			path: iss.File,
			lines: lines{
				begin: iss.Line,
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
